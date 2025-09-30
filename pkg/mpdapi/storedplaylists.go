package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
)

type StoredPlaylists interface {
	GetPlaylists() ([]Playlist, error)
	DeleteStoredPlaylist(string) error
	RenameStoredPlaylist(string, string) error
	SaveCurrentPlaylistAsStored(string) error
}

func (api *Impl) GetPlaylists() ([]Playlist, error) {
	cmd := commands.NewSingleCommand(commands.LISTPLAYLISTS)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	playlists, err := parser.ParseMultiValue[Playlist](list)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (api *Impl) DeleteStoredPlaylist(name string) error {
	cmd := commands.NewSingleCommand(commands.RM).AddParams(name)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) RenameStoredPlaylist(oldName, newName string) error {
	cmd := commands.NewSingleCommand(commands.RENAME).AddParams(oldName, newName)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) SaveCurrentPlaylistAsStored(name string) error {
	cmd := commands.NewSingleCommand(commands.SAVE).AddParams(name)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}
