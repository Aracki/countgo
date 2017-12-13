package yt

import (
	"context"
	"fmt"

	"github.com/aracki/countgo/yt/client"
	"google.golang.org/api/youtube/v3"
)

func InitYoutubeService() (*youtube.Service, error) {
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
		return nil, err
	}

	return s, err
}
