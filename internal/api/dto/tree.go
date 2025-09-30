package dto

type TreeItem interface{}

type DirectoryItem struct {
	Path     string     `json:"path"`
	Name     string     `json:"name"`
	Children []TreeItem `json:"children"`
}

type FileItem struct {
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	Time        *string `json:"time,omitempty"`
	Artist      *string `json:"artist,omitempty"`
	AlbumArtist *string `json:"album_artist,omitempty"`
	Title       *string `json:"title,omitempty"`
	Album       *string `json:"album,omitempty"`
	Track       *string `json:"track,omitempty"`
	Date        *string `json:"date,omitempty"`
}
