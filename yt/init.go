package yt

import (
	"context"
	"fmt"

	"github.com/aracki/countgo/yt/client"
	"google.golang.org/api/youtube/v3"
)

// InitYoutubeService will read from config file
// than make youtube client and service according to that client
// returns pointer to youtube.Service
func InitYoutubeService() (*youtube.Service, error) {
	ctx := context.Background()

	// reads from config file
	config, err := client.ReadConfigFile()
	if err != nil {
		fmt.Println("Unable to read/parse client secret file", err)
	}

	// making new client
	c := client.GetClient(ctx, config)

	// making new service based on client
	s, err := youtube.New(c)
	if err != nil {
		fmt.Println("Cannot make youtube client", err)
		return nil, err
	}

	return s, err
}
