package main

import (
	"github.com/aracki/countgo/youtube"
	"github.com/aracki/countgo/youtube/service"
)

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
)

func main() {

	// init youtube service
	s, _ := youtube.InitYoutubeService()

	// getting IvannSerbia channel info
	service.ChannelInfo(s, snippetContentDetailsStatistics, "IvannSerbia")

	// getting all the lists
	lists, _ := service.AllPlaylists(s, snippetContentDetails)
	// getting all the lists info concurrently
	service.PlaylistsInfo(s, snippetContentDetails, lists)
	// getting all the videos of all playlists of mine
	service.AllVideos(s, snippetContentDetails)
}
