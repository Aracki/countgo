package main

import (
	"strconv"

	"google.golang.org/api/youtube/v3"
)

// The appendPlaylistInfo goes through all videos in playlist (max 50)
// it uses response.NextPageToken to go to next 50 videos
// it populates *plInfoArr with new playlist info
func appendPlaylistInfo(service *youtube.Service, part string, playlist *youtube.Playlist, plInfoArr *[]string) {

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
}
