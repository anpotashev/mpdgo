package commands

import (
	"fmt"
)

type CommandType uint64

const (
	PLAY CommandType = iota
	PAUSE
	STOP
	PREV
	NEXT
	PLAYLIST_INFO
	STATUS
	LSINFO
	IDLE
	PING
	ENABLE_OUTPUT
	DISABLE_OUTPUT
	OUTPUTS
	CLEAR
	DELETE
	MOVE
	SHUFFLE
	ADD
	ADD_ID
	PLAY_ID
	SEEK
	LISTALL
	LISTALLINFO
	UPDATE
	LISTPLAYLISTS
	LISTPLAYLIST_INFO
	RANDOM
	REPEAT
	SINGLE
	CONSUME
	LOAD
	RM
	SAVE
	RENAME
	PASSWORD
)

func (c CommandType) String() string {
	switch c {
	case PLAY:
		return "play"
	case PAUSE:
		return "pause"
	case STOP:
		return "stop"
	case PREV:
		return "previous"
	case NEXT:
		return "next"
	case PLAYLIST_INFO:
		return "playlistinfo"
	case STATUS:
		return "status"
	case LSINFO:
		return "lsinfo"
	case IDLE:
		return "idle"
	case PING:
		return "ping"
	case ENABLE_OUTPUT:
		return "enableoutput"
	case DISABLE_OUTPUT:
		return "disableoutput"
	case OUTPUTS:
		return "outputs"
	case CLEAR:
		return "clear"
	case DELETE:
		return "delete"
	case MOVE:
		return "move"
	case SHUFFLE:
		return "shuffle"
	case ADD:
		return "add"
	case ADD_ID:
		return "addid"
	case PLAY_ID:
		return "playid"
	case SEEK:
		return "seek"
	case LISTALL:
		return "listall"
	case LISTALLINFO:
		return "listallinfo"
	case UPDATE:
		return "update"
	case LISTPLAYLISTS:
		return "listplaylists"
	case LISTPLAYLIST_INFO:
		return "listplaylistinfo"
	case RANDOM:
		return "random"
	case REPEAT:
		return "repeat"
	case SINGLE:
		return "single"
	case CONSUME:
		return "consume"
	case LOAD:
		return "load"
	case RM:
		return "rm"
	case SAVE:
		return "save"
	case RENAME:
		return "rename"
	case PASSWORD:
		return "password"
	default:
		return "unknown"
	}
}

type MpdCommand interface {
	fmt.Stringer
}
