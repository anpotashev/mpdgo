// package dto.Generated code. DO NOT EDIT.
package dto

import "github.com/anpotashev/mpdgo/pkg/mpdapi"

func MapOutput(in mpdapi.Output) Output {
	return Output{
		Id:      in.Id,
		Name:    in.Name,
		Enabled: in.Enabled,
	}
}

func MapSongTime(in mpdapi.SongTime) SongTime {
	return SongTime{
		Current: in.Current,
		Full:    in.Full,
	}
}

func MapStatus(in mpdapi.Status) Status {
	return Status{
		Volume:         in.Volume,
		Repeat:         in.Repeat,
		Random:         in.Random,
		Single:         in.Single,
		Consume:        in.Consume,
		Playlist:       in.Playlist,
		PlaylistLength: in.PlaylistLength,
		Xfade:          in.Xfade,
		State:          in.State,
		Song:           in.Song,
		SongId:         in.SongId,
		Bitrate:        in.Bitrate,
		Audio:          in.Audio,
		NextSong:       in.NextSong,
		NextSongId:     in.NextSongId,
		Time: func(in *mpdapi.SongTime) *SongTime {
			if in == nil {
				return nil
			}
			out := MapSongTime(*in)
			return &out
		}(in.Time),
	}
}

func MapPlaylist(in mpdapi.Playlist) Playlist {
	return Playlist{
		Name: in.Name,
		Items: func(items []mpdapi.PlaylistItem) (result []PlaylistItem) {
			for _, item := range items {
				result = append(result, MapPlaylistItem(item))
			}
			return result
		}(in.Items),
	}
}

func MapPlaylistItem(in mpdapi.PlaylistItem) PlaylistItem {
	return PlaylistItem{
		File:   in.File,
		Time:   in.Time,
		Artist: in.Artist,
		Title:  in.Title,
		Album:  in.Album,
		Track:  in.Track,
		Pos:    in.Pos,
		Id:     in.Id,
	}
}

func MapDirectoryItem(in mpdapi.DirectoryItem) DirectoryItem {
	return DirectoryItem{
		Path: in.Path,
		Name: in.Name,
	}
}

func MapFileItem(in mpdapi.FileItem) FileItem {
	return FileItem{
		Path:        in.Path,
		Name:        in.Name,
		Time:        in.Time,
		Artist:      in.Artist,
		AlbumArtist: in.AlbumArtist,
		Title:       in.Title,
		Album:       in.Album,
		Track:       in.Track,
		Date:        in.Date,
	}
}
