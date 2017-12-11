package service

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ResourceId  string `json:"title"`
	Thumbnail   string `json:"title"`
}

type Playlist struct {
	Title       string `json:"title"`
	VideosCount string `json:"videos_count"`
}
