package main

import (
	"context"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("/Users/raca/go/src/github.com/aracki/countgo/youtube/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	channelsListByUsername(service, "snippet,contentDetails,statistics", "IvannSerbia")
	l := getAllPlaylists(service, "snippet,contentDetails")
	videosAllPlaylists(service, "snippet,contentDetails", l)
}
