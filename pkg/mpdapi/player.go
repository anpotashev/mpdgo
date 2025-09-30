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

func (api *Impl) Play() error {
	cmd := commands.NewSingleCommand(commands.PLAY)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Pause() error {
	cmd := commands.NewSingleCommand(commands.PAUSE)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Stop() error {
	cmd := commands.NewSingleCommand(commands.STOP)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Previous() error {
	cmd := commands.NewSingleCommand(commands.PREV)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Next() error {
	cmd := commands.NewSingleCommand(commands.NEXT)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) PlayId(id int) error {
	cmd := commands.NewSingleCommand(commands.PLAY_ID).AddParams(id)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) PlayPos(pos int) error {
	cmd := commands.NewSingleCommand(commands.PLAY).AddParams(pos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) Seek(songPos, seekPos int) error {
	cmd := commands.NewSingleCommand(commands.SEEK).AddParams(songPos).AddParams(seekPos)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}
