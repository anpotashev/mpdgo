package dto

type SongTime struct {
	Current int `json:"current"`
	Full    int `json:"full"`
}

type Status struct {
	Volume         *int      `json:"volume,omitempty"`
	Repeat         *bool     `json:"repeat,omitempty"`
	Random         *bool     `json:"random,omitempty"`
	Single         *bool     `json:"single,omitempty"`
	Consume        *bool     `json:"consume,omitempty"`
	Playlist       *string   `json:"playlist,omitempty"`
	PlaylistLength *int      `json:"playlist_length,omitempty"`
	Xfade          *int      `json:"xfade,omitempty"`
	State          *string   `json:"state,omitempty"`
	Song           *int      `json:"song,omitempty"`
	SongId         *int      `json:"song_id,omitempty"`
	Time           *SongTime `json:"time,omitempty"`
	Bitrate        *int      `json:"bitrate,omitempty"`
	Audio          *string   `json:"audio,omitempty"`
	NextSong       *int      `json:"next_song,omitempty"`
	NextSongId     *int      `json:"next_song_id,omitempty"`
}
