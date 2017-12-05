package main

import (
	"context"
	"fmt"

	"google.golang.org/api/youtube/v3"
)

// getChannelInfo uses forUsername
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

// getAllPlaylists returns all playlist for current user
func getAllPlaylists(service *youtube.Service, part string) (playlists []*youtube.Playlist) {

	call := service.Playlists.List(part)
	// default maxResults is 5
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	handleError(err, "")

	var lists []*youtube.Playlist
	for _, item := range response.Items {
		lists = append(lists, item)
	}
	return lists
}

// showPlaylistInfo goes through all playlists
// and return videos count for each
func showPlaylistsInfo(service *youtube.Service, part string, playlists []*youtube.Playlist) {

	for _, playlist := range playlists {
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
				fmt.Println(playlist.Snippet.Title, ": ", pCount)
				break
			}
		}
	}
}

func main() {
	ctx := context.Background()

	config := readConfigFile()

	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	getChannelInfo(service, "snippet,contentDetails,statistics", "IvannSerbia")
	lists := getAllPlaylists(service, "snippet,contentDetails")
	showPlaylistsInfo(service, "snippet,contentDetails", lists)
}
