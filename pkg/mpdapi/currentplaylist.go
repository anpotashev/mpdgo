package mpdapi

import (
	"fmt"
	"time"

	"github.com/anpotashev/mpdgo/internal/commands"
	log "github.com/anpotashev/mpdgo/internal/logger"
	"github.com/anpotashev/mpdgo/internal/parser"
)

type Playlist struct {
	Items        []PlaylistItem
	Name         *string    `mpd_prefix:"playlist" is_new_element_prefix:"true"`
	LastModified *time.Time `mpd_prefix:"Last-Modified"`
}

type PlaylistItem struct {
	File   string  `mpd_prefix:"file" is_new_element_prefix:"true"`
	Time   int     `mpd_prefix:"Time"`
	Artist *string `mpd_prefix:"Artist"`
	Title  *string `mpd_prefix:"Title"`
	Album  *string `mpd_prefix:"Album"`
	Track  *string `mpd_prefix:"Track"`
	Pos    int     `mpd_prefix:"Pos"`
	Id     int     `mpd_prefix:"Id"`
}

type CurrentPlaylist interface {
	Playlist() (*Playlist, error)
	PlaylistInfo(name string) (*Playlist, error)
	Clear() error
	Add(path string) error
	AddToPos(pos int, path string) error
	DeleteByPos(pos int) error
	Move(fromPos, toPos int) error
	BatchMove(fromStartPos, fromEndPos, toPos int) error
	ShuffleAll() error
	Shuffle(fromPos, toPos int) error
	AddStoredToPos(name string, pos int) error
}

func (api *Impl) Playlist() (*Playlist, error) {
	cmd := commands.NewSingleCommand(commands.PLAYLIST_INFO)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
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
func (api *Impl) PlaylistInfo(name string) (*Playlist, error) {
	cmd := commands.NewSingleCommand(commands.LISTPLAYLIST_INFO).AddParams(name)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
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
	cmd := commands.NewSingleCommand(commands.CLEAR)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Add(path string) error {
	log.DebugContext(api.requestContext, "getting files by path")
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		log.DebugContext(api.requestContext, "paths is empty")
		return nil
	}
	log.DebugContext(api.requestContext, "paths is not empty. Making batch-command")
	var cmds []commands.SingleCommand
	for _, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD).AddParams(p)
		cmds = append(cmds, cmd)
	}
	return wrapPkgError(api.mpdClient.SendBatchCommand(api.requestContext, cmds))
}

func (api *Impl) AddToPos(pos int, path string) error {
	paths, err := api.getFilesPaths(path)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return nil
	}
	var cmds []commands.SingleCommand
	for i, p := range paths {
		cmd := commands.NewSingleCommand(commands.ADD_ID).AddParams(p).AddParams(i + pos)
		cmds = append(cmds, cmd)
	}
	return wrapPkgError(api.mpdClient.SendBatchCommand(api.requestContext, cmds))
}

func (api *Impl) getFilesPaths(path string) ([]string, error) {
	cmd := commands.NewSingleCommand(commands.LISTALL).AddParams(path)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	//log.DebugContext(api.requestContext, "list", list)
	items, err := parser.ParseMultiValue[ParsedItem](list)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	log.DebugContext(api.requestContext, "after parsing", "items size", len(items))
	var paths []string
	for _, item := range items {
		log.DebugContext(api.requestContext, "processing", "item", item)
		if item.File != nil {
			log.DebugContext(api.requestContext, "added")
			paths = append(paths, *item.File)
		}
	}
	log.DebugContext(api.requestContext, "after parsing", "paths size", len(paths))
	return paths, nil
}

func (api *Impl) DeleteByPos(pos int) error {
	cmd := commands.NewSingleCommand(commands.DELETE).AddParams(pos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Move(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE).AddParams(fromPos).AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) BatchMove(fromStartPos, fromEndPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.MOVE).AddParams(fmt.Sprintf("%d:%d", fromStartPos, fromEndPos)).AddParams(toPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) ShuffleAll() error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Shuffle(fromPos, toPos int) error {
	cmd := commands.NewSingleCommand(commands.SHUFFLE).AddParams(fmt.Sprintf("%d:%d", fromPos, toPos))
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) AddStoredToPos(name string, pos int) error {
	cmd := commands.NewSingleCommand(commands.LISTPLAYLIST_INFO).AddParams(name)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
	if err != nil {
		return wrapPkgError(err)
	}
	playlistItems, err := parser.ParseMultiValue[PlaylistItem](list)
	if err != nil {
		return wrapPkgError(err)
	}
	var cmds []commands.SingleCommand
	for i, item := range playlistItems {
		cmds = append(cmds, commands.NewSingleCommand(commands.ADD_ID).AddParams(item.File).AddParams(pos+i))
	}
	return wrapPkgError(api.mpdClient.SendBatchCommand(api.requestContext, cmds))
}
