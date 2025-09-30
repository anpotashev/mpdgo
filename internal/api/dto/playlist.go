package dto

import "time"

type Playlist struct {
	Items        []PlaylistItem `json:"items,omitempty"`
	Name         *string        `json:"name"`
	LastModified *time.Time     `json:"last_modified,omitempty"`
}

type PlaylistItem struct {
	File   string  `json:"file"`
	Time   int     `json:"time"`
	Artist *string `json:"artist,omitempty"`
	Title  *string `json:"title,omitempty"`
	Album  *string `json:"album,omitempty"`
	Track  *string `json:"track,omitempty"`
	Pos    int     `json:"pos"`
	Id     int     `json:"id"`
}
