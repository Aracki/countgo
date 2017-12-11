package service

import (
	"fmt"
	"sync"

	"google.golang.org/api/youtube/v3"
)

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
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

// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func ChannelInfo(service *youtube.Service, forUsername string) (string, error) {

	call := service.Channels.List(snippetContentDetailsStatistics)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	var info string

	info = fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views. \n",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount)
	info += fmt.Sprintf(response.Items[0].Snippet.Description)

	return info, nil
}

// The AllPlaylists uses current user - maxResult is set to 50 (default is 5)
// It runs go routines for each playlist
// and call appendPlaylistInfo which populates plInfo array.
// Different goroutines are appending the same slice;
// WaitGroup waits for all goroutines to finish
func AllPlaylists(service *youtube.Service) ([]Playlist, error) {

	// get all playlists
	call := service.Playlists.List(snippetContentDetails)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(response.Items))

	var pls []Playlist
	for _, pl := range response.Items {
		go func(p *youtube.Playlist) {
			appendPlaylistInfo(service, snippetContentDetails, p, &pls)
			wg.Done()
		}(pl)
	}
	wg.Wait()

	return pls, nil
}

// Gets all the videos of specific youtube.Playlist
func AllVideosByPlaylist(service *youtube.Service, pl *youtube.Playlist) ([]Video, error) {

	var vds []Video
	pageToken := ""

	for {
		call := service.PlaylistItems.List(snippetContentDetails)
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
func AllVideos(service *youtube.Service) ([]Video, error) {

	var vds []Video

	// get all playlists of mine
	call := service.Playlists.List(snippetContentDetails)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(response.Items))

	for _, pl := range response.Items {
		go func(p *youtube.Playlist) {
			v, _ := AllVideosByPlaylist(service, p)
			vds = append(vds, v...)
			wg.Done()
		}(pl)
	}
	wg.Wait()

	return vds, nil
}
