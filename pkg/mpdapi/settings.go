package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
	"regexp"
	"strconv"
)

type Settings interface {
	Random(value bool) error
	Repeat(value bool) error
	Single(value bool) error
	Consume(value bool) error
	Status() (Status, error)
}

type SongTime struct {
	Current int
	Full    int
}
type status struct {
	Volume         int    `mpd_prefix:"volume"`
	Repeat         bool   `mpd_prefix:"repeat"`
	Random         bool   `mpd_prefix:"random"`
	Single         bool   `mpd_prefix:"single"`
	Consume        bool   `mpd_prefix:"consume"`
	Playlist       string `mpd_prefix:"playlist"`
	PlaylistLength int    `mpd_prefix:"playlistlength"`
	Xfade          int    `mpd_prefix:"xfade"`
	State          string `mpd_prefix:"state"`
	Song           int    `mpd_prefix:"song"`
	SongId         int    `mpd_prefix:"songid"`
	Time           string `mpd_prefix:"time"`
	Bitrate        int    `mpd_prefix:"bitrate"`
	Audio          string `mpd_prefix:"audio"`
	NextSong       int    `mpd_prefix:"nextsong"`
	NextSongId     int    `mpd_prefix:"nextsongid"`
}

type Status struct {
	Volume         int
	Repeat         bool
	Random         bool
	Single         bool
	Consume        bool
	Playlist       string
	PlaylistLength int
	Xfade          int
	State          string
	Song           int
	SongId         int
	Time           SongTime
	Bitrate        int
	Audio          string
	NextSong       int
	NextSongId     int
}

func (api *MpdApiImpl) Random(value bool) error {
	cmd := commands.NewSingleCommand(commands.RANDOM)
	cmd, _ = cmd.AddParams(value)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Repeat(value bool) error {
	cmd := commands.NewSingleCommand(commands.REPEAT)
	cmd, _ = cmd.AddParams(value)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Single(value bool) error {
	cmd := commands.NewSingleCommand(commands.SINGLE)
	cmd, _ = cmd.AddParams(value)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}
func (api *MpdApiImpl) Consume(value bool) error {
	cmd := commands.NewSingleCommand(commands.CONSUME)
	cmd, _ = cmd.AddParams(value)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Status() (Status, error) {
	cmd := commands.NewSingleCommand(commands.STATUS)
	list, err := api.mpdClient.SendCommand(cmd)
	if err != nil {
		return Status{}, wrapPkgError(err)
	}
	status, err := parser.ParseSingleValue[status](list)
	if err != nil {
		return Status{}, wrapPkgError(err)
	}
	songTime := SongTime{}
	songTimeRegexp := regexp.MustCompile("(\\d+):(\\d+)")
	matches := songTimeRegexp.FindStringSubmatch(status.Time)
	if len(matches) == 3 {
		songTime.Full, _ = strconv.Atoi(matches[1])
		songTime.Current, _ = strconv.Atoi(matches[2])
	}
	result := Status{
		Volume:         status.Volume,
		Repeat:         status.Repeat,
		Random:         status.Random,
		Single:         status.Single,
		Consume:        status.Consume,
		Playlist:       status.Playlist,
		PlaylistLength: status.PlaylistLength,
		Xfade:          status.Xfade,
		State:          status.State,
		Song:           status.Song,
		SongId:         status.SongId,
		Time:           songTime,
		Bitrate:        status.Bitrate,
		Audio:          status.Audio,
		NextSong:       status.NextSong,
		NextSongId:     status.NextSongId,
	}
	return result, nil
}
