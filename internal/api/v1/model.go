package v1

// current playlist
type AddToPlaylistRequest struct {
	Path string `json:"path"`
}
type AddToPosPlaylistRequest struct {
	Path string `json:"path"`
	Pos  int    `json:"pos"`
}
type DeleteByPosFromCurrentPlaylistRequest struct {
	Pos int `json:"pos"`
}

type ShuffleCurrentPlaylistRequest struct {
	FromPos int `json:"from_pos"`
	ToPos   int `json:"to_pos"`
}
type MoveInCurrentPlaylistRequest struct {
	FromPos int `json:"from_pos"`
	ToPos   int `json:"to_pos"`
}

type BatchMoveInCurrentPlaylistRequest struct {
	FromStartPos int `json:"from_start_pos"`
	FromEndPos   int `json:"from_end_pos"`
	ToPos        int `json:"to_pos"`
}

// output
type ChangeOutputStateRequest struct {
	OutputId int `json:"output_id"`
}

type PlayIdRequest struct {
	Id int `json:"id"`
}

type PlayPosRequest struct {
	Pos int `json:"pos"`
}

type SeekPosRequest struct {
	Pos      int `json:"pos"`
	SeekTime int `json:"seek_time"`
}

type GetStoredPlaylistsRequest struct {
	Name string `json:"name"`
}

type DeleteStoredPlaylistRequest struct {
	Name string `json:"name"`
}

type SaveCurrentPlaylistAsStoredRequest struct {
	Name string `json:"name"`
}

type RenameStoredPlaylistRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}
