package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"google.golang.org/api/youtube/v3"
)

// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func getChannelInfo(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views. \n",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
	fmt.Println(response.Items[0].Snippet.Description, "\n")
}

// The getAllPlaylists uses current user
// maxResult is set to 50 (default is 5)
// returns all playlists
func getAllPlaylists(service *youtube.Service, part string) (playlists []*youtube.Playlist) {

	call := service.Playlists.List(part)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	handleError(err, "")

	var lists []*youtube.Playlist
	for _, item := range response.Items {
		lists = append(lists, item)
	}
	return lists
}

// The getPlaylistsInfo runs go routines for each playlist
// and call appendPlaylistInfo which populates plInfo array.
// Different goroutines are appending the same slice,
// WaitGroup waits for all goroutines to finish
func getPlaylistsInfo(service *youtube.Service, part string, playlists []*youtube.Playlist) {

	var wg sync.WaitGroup
	wg.Add(len(playlists))

	var plInfoArr []string
	for _, playlist := range playlists {
		go appendPlaylistInfo(service, part, playlist, &plInfoArr, &wg)
	}
	wg.Wait()

	fmt.Println(plInfoArr)
}

// The appendPlaylistInfo goes through all videos in playlist (max 50)
// it uses response.NextPageToken to go to next 50 videos
// it populates *plInfoArr with new playlist info
func appendPlaylistInfo(service *youtube.Service, part string, playlist *youtube.Playlist, plInfoArr *[]string, wg *sync.WaitGroup) {

	pageToken := ""
	pCount := 0
	for {
		call := service.PlaylistItems.List(part)
		call = call.PlaylistId(playlist.Id).MaxResults(50)
		response, err := call.PageToken(pageToken).Do()
		handleError(err, "")

		// increment counter and move to another page of 50 videos
		pCount += len(response.Items)
		pageToken = response.NextPageToken

		if pageToken == "" {
			*plInfoArr = append(*plInfoArr, playlist.Snippet.Title, ":[", strconv.Itoa(pCount), "]")
			break
		}
	}
	wg.Done()
}

func main() {
	ctx := context.Background()

	config := readConfigFile()

	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	getChannelInfo(service, "snippet,contentDetails,statistics", "IvannSerbia")
	lists := getAllPlaylists(service, "snippet,contentDetails")
	getPlaylistsInfo(service, "snippet,contentDetails", lists)
}
