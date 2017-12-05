package main

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// readConfigFile will return oauth2 config
// based on client_secret.json which is located in project root
func readConfigFile() *oauth2.Config {

	filePath, _ := filepath.Abs("../client_secret.json")

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config
}

func main() {
	ctx := context.Background()

	config := readConfigFile()

	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	channelsListByUsername(service, "snippet,contentDetails,statistics", "IvannSerbia")
	l := getAllPlaylists(service, "snippet,contentDetails")
	videosAllPlaylists(service, "snippet,contentDetails", l)
}
