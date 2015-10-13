package main

import (
	"github.com/Sirupsen/logrus"
)

type TrimCommand struct {
	options *Options
	stream  *logrus.Logger
}

func (c TrimCommand) Execute(args []string) error {
	c.options.Init(c.stream)

	return nil
}

func NewTrimCommand(stream *logrus.Logger, options *Options) *TrimCommand {
	return &TrimCommand{
		options: options,
		stream:  stream,
	}
}
