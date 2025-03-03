package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
)

type StoredPlaylists interface {
	GetPlaylists() ([]Playlist, error)
}

func (api *Impl) GetPlaylists() ([]Playlist, error) {
	cmd := commands.NewSingleCommand(commands.LISTPLAYLISTS)
	list, err := api.mpdClient.SendCommand(cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	playlists, err := parser.ParseMultiValue[Playlist](list)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}
