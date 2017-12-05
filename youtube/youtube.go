package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `Please configure OAuth 2.0`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// channelsListByUsername uses forUsername
// to get Channel info (id, tittle, totalViews and description)
func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
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

// showPlaylistInfo uses playlist
// to return all videos count
func showPlaylistInfo(service *youtube.Service, part string, playlist *youtube.Playlist) {

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

// for loop through all playlists
func videosAllPlaylists(service *youtube.Service, part string, playlists []*youtube.Playlist) {

	for _, playlist := range playlists {
		showPlaylistInfo(service, part, playlist)
	}
}
