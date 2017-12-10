package main

import (
	"context"
	"fmt"

	"github.com/aracki/countgo/youtube/client"
	"github.com/aracki/countgo/youtube/service"
	"google.golang.org/api/youtube/v3"
)

const (
	snippetContentDetailsStatistics = "snippet,contentDetails,statistics"
	snippetContentDetails           = "snippet,contentDetails"
)

func main() {
	ctx := context.Background()

	// reads from config file
	config, err := client.ReadConfigFile()
	if err != nil {
		fmt.Println("Unable to read/parse client secret file", err)
	}

	// making new client
	c := client.GetClient(ctx, config)
	s, err := youtube.New(c)
	if err != nil {
		fmt.Println("Cannot make youtube client", err)
	}

	// getting IvannSerbia channel info
	service.ChannelInfo(s, snippetContentDetailsStatistics, "IvannSerbia")

	// getting all the lists
	lists, _ := service.AllPlaylists(s, snippetContentDetails)
	// getting all the lists info concurrently
	service.PlaylistsInfo(s, snippetContentDetails, lists)
	// getting all the videos of all playlists of mine
	service.AllVideos(s, snippetContentDetails)
}
