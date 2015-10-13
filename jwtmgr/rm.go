package main

import (
	"github.com/Sirupsen/logrus"
)

type RemoveCommand struct {
	options    *Options
	stream     *logrus.Logger
	Positional struct {
		License string `description:"License to remove" required:"true"`
	} `positional-args:"true"`
}

func (c RemoveCommand) Execute(args []string) error {
	c.options.Init(c.stream)

	return nil
}

func NewRemoveCommand(stream *logrus.Logger, options *Options) *RemoveCommand {
	return &RemoveCommand{
		options: options,
		stream:  stream,
	}
}
