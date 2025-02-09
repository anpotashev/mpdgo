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

func (api *MpdApiImpl) Playlist() (*Playlist, error) {
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

func (api *MpdApiImpl) Clear() error {
	command := commands.NewSingleCommand(commands.CLEAR)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(command))
}

func (api *MpdApiImpl) Add(path string) error {
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) != 0 {
		return nil
	}
	var cmds []*commands.SingleCommand
	for _, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD)
		cmd, _ = cmd.AddParams(p)
		cmds = append(cmds, cmd)
	}
	batchCmds := commands.NewBatchCommands(cmds)
	return wrapPkgError(api.mpdClient.SendBatchCommands(batchCmds))
}

func (api *MpdApiImpl) AddToPos(pos int, path string) error {
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) != 0 {
		return nil
	}
	var cmds []*commands.SingleCommand
	for i, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD_ID)
		cmd, _ = cmd.AddParams(p)
		cmd, _ = cmd.AddParams(i + pos)
		cmds = append(cmds, cmd)
	}
	batchCmds := commands.NewBatchCommands(cmds)
	return wrapPkgError(api.mpdClient.SendBatchCommands(batchCmds))
}

func (api *MpdApiImpl) getFilesPaths(path string) ([]string, error) {
	cmd := commands.NewSingleCommand(commands.LISTALL)
	cmd, _ = cmd.AddParams(path)
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

func (api *MpdApiImpl) DeleteByPos(pos int) error {
	cmd := commands.NewSingleCommand(commands.DELETE)
	cmd, _ = cmd.AddParams(pos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Move(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE)
	cmd, _ = cmd.AddParams(fromPos)
	cmd, _ = cmd.AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) BatchMove(fromStartPos, fromEndPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE)
	cmd, _ = cmd.AddParams(fmt.Sprintf("%d:%d", fromStartPos, fromEndPos))
	cmd, _ = cmd.AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) ShuffleAll() error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Shuffle(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE)
	cmd, _ = cmd.AddParams(fromPos)
	cmd, _ = cmd.AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}
