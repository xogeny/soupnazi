package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/xogeny/soupnazi"
)

type ListCommand struct {
	options *Options
	stream  *logrus.Logger
	Raw     bool `long:"raw" short:"r" description:"Dump licenses as raw test"`
}

func (c ListCommand) Execute(args []string) error {
	c.options.Init(c.stream)

	lfile := soupnazi.LicenseFile()

	licenses, err := soupnazi.ParseLicenses(lfile)
	if err != nil {
		return err
	}

	for _, s := range licenses {
		fmt.Printf("%s\n", s)
	}

	return nil
}

func NewListCommand(stream *logrus.Logger, options *Options) *ListCommand {
	return &ListCommand{
		options: options,
		stream:  stream,
	}
}
