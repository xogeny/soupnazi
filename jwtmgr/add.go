package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/xogeny/soupnazi"
)

type AddCommand struct {
	options    *Options
	stream     *logrus.Logger
	Positional struct {
		License string `description:"License to add" required:"true"`
	} `positional-args:"true"`
}

func (c AddCommand) Execute(args []string) error {
	c.options.Init(c.stream)

	if c.Positional.License == "" {
		return fmt.Errorf("No license specified")
	}

	err := soupnazi.AddLicense(c.Positional.License, c.stream)
	if err != nil {
		return err
	}

	fmt.Printf("License successfully installed\n")
	return nil
}

func NewAddCommand(stream *logrus.Logger, options *Options) *AddCommand {
	return &AddCommand{
		options: options,
		stream:  stream,
	}
}
