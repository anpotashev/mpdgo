package mpdapi

import (
	"fmt"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
	"time"
)

type Playlist struct {
	Items        []PlaylistItem
	Name         string    `mpd_prefix:""`
	LastModified time.Time `mpd_prefix:"Last-Modified"`
}

type PlaylistItem struct {
	File   string `mpd_prefix:"file" is_new_element_prefix:"true"`
	Time   int    `mpd_prefix:"Time"`
	Artist string `mpd_prefix:"Artist"`
	Title  string `mpd_prefix:"Title"`
	Album  string `mpd_prefix:"Album"`
	Track  string `mpd_prefix:"Track"`
	Pos    int    `mpd_prefix:"Pos"`
	Id     int    `mpd_prefix:"Id"`
}

type CurrentPlaylist interface {
	Playlist() (*Playlist, error)
	Clear() error
	Add(path string) error
	AddToPos(pos int, path string) error
	DeleteByPos(pos int) error
	Move(fromPos, toPos int) error
	BatchMove(fromStartPos, fromEndPos, toPos int) error
	ShuffleAll() error
	Shuffle(fromPos, toPos int) error
}

func (api *Impl) Playlist() (*Playlist, error) {
	command := commands.NewSingleCommand(commands.PLAYLIST_INFO)
	list, err := api.mpdClient.SendCommand(command)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	playlistItems, err := parser.ParseMultiValue[PlaylistItem](list)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	return &Playlist{
		Items: playlistItems,
	}, nil
}

func (api *Impl) Clear() error {
	command := commands.NewSingleCommand(commands.CLEAR)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(command))
}

func (api *Impl) Add(path string) error {
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) != 0 {
		return nil
	}
	var cmds []*commands.SingleCommand
	for _, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD).AddParams(p)
		cmds = append(cmds, cmd)
	}
	batchCmds := commands.NewBatchCommands(cmds)
	return wrapPkgError(api.mpdClient.SendBatchCommands(batchCmds))
}

func (api *Impl) AddToPos(pos int, path string) error {
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) != 0 {
		return nil
	}
	var cmds []*commands.SingleCommand
	for i, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD_ID).AddParams(p).AddParams(i + pos)
		cmds = append(cmds, cmd)
	}
	batchCmds := commands.NewBatchCommands(cmds)
	return wrapPkgError(api.mpdClient.SendBatchCommands(batchCmds))
}

func (api *Impl) getFilesPaths(path string) ([]string, error) {
	cmd := commands.NewSingleCommand(commands.LISTALL).AddParams(path)
	list, err := api.mpdClient.SendCommand(cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	items, err := parser.ParseMultiValue[ParsedItem](list)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	var paths []string
	for _, item := range items {
		if item.File != "" {
			paths = append(paths, item.File)
		}
	}
	return paths, nil
}

func (api *Impl) DeleteByPos(pos int) error {
	cmd := commands.NewSingleCommand(commands.DELETE).AddParams(pos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *Impl) Move(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE).AddParams(fromPos).AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *Impl) BatchMove(fromStartPos, fromEndPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE).AddParams(fmt.Sprintf("%d:%d", fromStartPos, fromEndPos)).AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *Impl) ShuffleAll() error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *Impl) Shuffle(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE).AddParams(fromPos).AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}
