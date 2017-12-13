package main

import (
	"fmt"

	"github.com/aracki/countgo/youtube/service"
	"github.com/aracki/countgo/yt"
)

func main() {

	// init youtube service
	s, _ := yt.InitYoutubeService()

	// getting IvannSerbia channel info
	info, _ := service.ChannelInfo(s, "IvannSerbia")
	fmt.Println(info)

	// getting all the lists info concurrently
	lists, _ := service.AllPlaylists(s)
	for _, v := range lists {
		fmt.Printf("%+v", v)
	}
	// getting all the videos of all playlists of mine
	service.AllVideos(s)
}
