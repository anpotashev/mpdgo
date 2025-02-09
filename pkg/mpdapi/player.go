package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
)

type Player interface {
	Play() error
	Pause() error
	Stop() error
	Previous() error
	Next() error
	PlayId(id int) error
	PlayPos(pos int) error
	Seek(songPos, seekPos int) error
}

func (api *MpdApiImpl) Play() error {
	cmd := commands.NewSingleCommand(commands.PLAY)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Pause() error {
	cmd := commands.NewSingleCommand(commands.PAUSE)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Stop() error {
	cmd := commands.NewSingleCommand(commands.STOP)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Previous() error {
	cmd := commands.NewSingleCommand(commands.PREV)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Next() error {
	cmd := commands.NewSingleCommand(commands.NEXT)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) PlayId(id int) error {
	cmd := commands.NewSingleCommand(commands.PLAY)
	cmd, _ = cmd.AddParams(id)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) PlayPos(pos int) error {
	cmd := commands.NewSingleCommand(commands.PLAY_ID)
	cmd, _ = cmd.AddParams(pos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}

func (api *MpdApiImpl) Seek(songPos, seekPos int) error {
	cmd := commands.NewSingleCommand(commands.SEEK)
	cmd, _ = cmd.AddParams(songPos)
	cmd, _ = cmd.AddParams(seekPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendCommand(cmd))
}
