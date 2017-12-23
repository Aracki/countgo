package service

import (
	"strconv"

	"github.com/aracki/countgo/models"
	"google.golang.org/api/youtube/v3"
)

// The appendPlaylistInfo goes through all videos in playlist (max 50)
// it uses response.NextPageToken to go to next 50 videos
// it populates *plInfoArr with new playlist info
func appendPlaylistInfo(service *youtube.Service, part string, playlist *youtube.Playlist, plInfoArr *[]models.Playlist) error {

	pageToken := ""
	pCount := 0
	for {
		call := service.PlaylistItems.List(part)
		call = call.PlaylistId(playlist.Id).MaxResults(50)
		response, err := call.PageToken(pageToken).Do()
		if err != nil {
			return err
		}

		// increment counter and move to another page of 50 videos
		pCount += len(response.Items)
		pageToken = response.NextPageToken

		if pageToken == "" {
			// append total count to plInfoArr
			*plInfoArr = append(*plInfoArr, models.Playlist{
				Id:          playlist.Id,
				Title:       playlist.Snippet.Title,
				VideosCount: strconv.Itoa(pCount),
			})
			break
		}
	}
	return nil
}
