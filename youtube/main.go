package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/aracki/countgo/youtube/client"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ResourceId  string `json:"title"`
	Thumbnail   string `json:"title"`
}

type Playlist struct {
	Title       string `json:"title"`
	VideosCount string `json:"videos_count"`
}

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
)

// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func ChannelInfo(service *youtube.Service, part string, forUsername string) error {

	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views. \n",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
	fmt.Println(response.Items[0].Snippet.Description, "\n")

	return nil
}

// The AllPlaylists uses current user
// maxResult is set to 50 (default is 5)
// returns all playlists
func AllPlaylists(service *youtube.Service, part string) ([]*youtube.Playlist, error) {

	call := service.Playlists.List(part)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var lists []*youtube.Playlist
	for _, item := range response.Items {
		lists = append(lists, item)
	}
	return lists, nil
}

// The PlaylistsInfo runs go routines for each playlist
// and call appendPlaylistInfo which populates plInfo array.
// Different goroutines are appending the same slice,
// WaitGroup waits for all goroutines to finish
func PlaylistsInfo(service *youtube.Service, part string, playlists []*youtube.Playlist) ([]Playlist, error) {

	var wg sync.WaitGroup
	wg.Add(len(playlists))

	var pls []Playlist
	for _, playlist := range playlists {
		go func(pl *youtube.Playlist) {
			appendPlaylistInfo(service, part, pl, &pls)
			wg.Done()
		}(playlist)
	}
	wg.Wait()

	return pls, nil
}

// Gets all the videos of specific youtube.Playlist
func AllVideosByPlaylist(service *youtube.Service, part string, pl *youtube.Playlist) ([]Video, error) {

	var vds []Video
	pageToken := ""

	for {
		call := service.PlaylistItems.List(part)
		call = call.PlaylistId(pl.Id).MaxResults(50)
		response, err := call.PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}

		// move pageToken to another page
		pageToken = response.NextPageToken

		for _, item := range response.Items {
			t := ""
			if item.Snippet.Thumbnails != nil && item.Snippet.Thumbnails.Medium != nil {
				t = item.Snippet.Thumbnails.Medium.Url
			}
			vds = append(vds, Video{
				item.Snippet.Title,
				item.Snippet.Description,
				item.Snippet.PublishedAt,
				item.Snippet.ResourceId.VideoId,
				t,
			})
		}
		// if there are no pages
		if pageToken == "" {
			break
		}
	}
	return vds, nil
}

// Gets all the videos of all playlists of mine
// goes through all playlists and concurrently appending to vds array of Videos
func AllVideos(service *youtube.Service, part string) ([]Video, error) {

	var vds []Video

	// get all playlists of mine
	call := service.Playlists.List(part)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(response.Items))

	for _, pl := range response.Items {
		go func(p *youtube.Playlist) {
			v, _ := AllVideosByPlaylist(service, part, p)
			vds = append(vds, v...)
			wg.Done()
		}(pl)
	}
	wg.Wait()

	return vds, nil
}

func main() {
	ctx := context.Background()

	// reads from config file
	config, err := client.ReadConfigFile()
	if err != nil {
		fmt.Println("Unable to read/parse client secret file", err)
	}

	// making new client
	c := client.GetClient(ctx, config)
	service, err := youtube.New(c)
	if err != nil {
		fmt.Println("Cannot make youtube client", err)
	}

	// getting IvannSerbia channel info
	ChannelInfo(service, snippetContentDetailsStatistics, "IvannSerbia")

	// getting all the lists
	lists, _ := AllPlaylists(service, snippetContentDetails)
	// getting all the lists info concurrently
	PlaylistsInfo(service, snippetContentDetails, lists)
	// getting all the videos of all playlists of mine
	AllVideos(service, snippetContentDetails)
}
