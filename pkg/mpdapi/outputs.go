package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/parser"
)

type Outputs interface {
	EnableOutput(id int) error
	DisableOutput(id int) error
	ListOutputs() ([]Output, error)
}

type Output struct {
	Name    string `mpd_prefix:"outputname"`
	Id      int    `mpd_prefix:"outputid" is_new_element_prefix:"true"`
	Enabled bool   `mpd_prefix:"outputenabled"`
}

func (api *Impl) EnableOutput(id int) error {
	cmd := commands.NewSingleCommand(commands.ENABLE_OUTPUT).AddParams(id)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) DisableOutput(id int) error {
	cmd := commands.NewSingleCommand(commands.DISABLE_OUTPUT).AddParams(id)
	return wrapPkgErrorIgnoringAnswer(api.mpdClient.SendSingleCommand(api.requestContext, cmd))
}

func (api *Impl) ListOutputs() ([]Output, error) {
	cmd := commands.NewSingleCommand(commands.OUTPUTS)
	list, err := api.mpdClient.SendSingleCommand(api.requestContext, cmd)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	result, err := parser.ParseMultiValue[Output](list)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	return result, nil
}
