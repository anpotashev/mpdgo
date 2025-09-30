package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
	"strings"
)

type Tree interface {
	Tree() (*DirectoryItem, error)
	UpdateDB(path string) error
}

type TreeItem interface {
	getParent() *DirectoryItem
	getName() string
	isLeaf() bool
}

type DirectoryItem struct {
	parent   *DirectoryItem
	Path     string
	Name     string
	Children []TreeItem
}

type FileItem struct {
	parent      *DirectoryItem
	Path        string
	Name        string
	Time        *string
	Artist      *string
	AlbumArtist *string
	Title       *string
	Album       *string
	Track       *string
	Date        *string
}

func (d *DirectoryItem) getParent() *DirectoryItem {
	return d.parent
}

func (d *DirectoryItem) getName() string {
	return d.Name
}

func (d *DirectoryItem) isLeaf() bool {
	return false
}

func (f *FileItem) getParent() *DirectoryItem {
	return f.parent
}

func (f *FileItem) getName() string {
	return f.Name
}

func (f *FileItem) isLeaf() bool {
	return true
}

type ParsedItem struct {
	File        *string `mpd_prefix:"file" is_new_element_prefix:"true"`
	Directory   *string `mpd_prefix:"directory" is_new_element_prefix:"true"`
	Time        *string `mpd_prefix:"Time"`
	Artist      *string `mpd_prefix:"Artist"`
	AlbumArtist *string `mpd_prefix:"AlbumArtist"`
	Title       *string `mpd_prefix:"Title"`
	Album       *string `mpd_prefix:"Album"`
	Track       *string `mpd_prefix:"Track"`
	Date        *string `mpd_prefix:"Date"`
	Genre       *string `mpd_prefix:"Genre"`
}

func (api *Impl) Tree() (*DirectoryItem, error) {
	cmd := commands.NewSingleCommand(commands.LISTALLINFO)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	mpdParsedItems, err := parser.ParseMultiValue[ParsedItem](list)
	if err != nil {
		return nil, err
	}
	rootItem := &DirectoryItem{
		parent:   nil,
		Name:     "/",
		Path:     "",
		Children: make([]TreeItem, 0),
	}
	currentDir := rootItem
	for _, item := range mpdParsedItems {
		if item.Directory != nil {
			parentDirItem := findParentDirItem(*item.Directory, currentDir)
			name := strings.TrimPrefix(*item.Directory, parentDirItem.Path)
			name = strings.TrimPrefix(name, "/")
			dirItem := &DirectoryItem{
				parent:   parentDirItem,
				Name:     name,
				Path:     *item.Directory,
				Children: make([]TreeItem, 0),
			}
			parentDirItem.Children = append(parentDirItem.Children, dirItem)
			currentDir = dirItem
		} else {
			parentDirItem := findParentDirItem(*item.File, currentDir)
			name := strings.TrimPrefix(*item.File, parentDirItem.Path)
			name = strings.TrimPrefix(name, "/")
			fileItem := &FileItem{
				parent:      parentDirItem,
				Name:        name,
				Path:        *item.File,
				Time:        item.Time,
				Artist:      item.Artist,
				AlbumArtist: item.AlbumArtist,
				Title:       item.Title,
				Album:       item.Album,
				Track:       item.Track,
				Date:        item.Date,
			}
			parentDirItem.Children = append(parentDirItem.Children, fileItem)
			currentDir = parentDirItem
		}
	}
	return rootItem, nil
}

func findParentDirItem(path string, currentActiveDir *DirectoryItem) *DirectoryItem {
	if strings.HasPrefix(path, currentActiveDir.Path) {
		return currentActiveDir
	}
	return findParentDirItem(path, currentActiveDir.parent)
}

func (api *Impl) UpdateDB(path string) error {
	cmd := commands.NewSingleCommand(commands.UPDATE).AddParams(path)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}
