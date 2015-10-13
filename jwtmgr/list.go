package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/xogeny/soupnazi"
)

type featureInfo struct {
}

type appInfo map[string]*featureInfo

func (ai appInfo) Feature(name string) *featureInfo {
	_, exists := ai[name]
	if !exists {
		ai[name] = &featureInfo{}
	}
	return ai[name]
}

type AppCollection map[string]*appInfo

func (ac AppCollection) App(name string) *appInfo {
	_, exists := ac[name]
	if !exists {
		ac[name] = &appInfo{}
	}
	return ac[name]
}

type ListCommand struct {
	options *Options
	stream  *logrus.Logger
	Raw     bool `long:"raw" short:"r" description:"Dump licenses as raw test"`
}

func (c ListCommand) Execute(args []string) error {
	c.options.Init(c.stream)

	lfile := soupnazi.LicenseFile()

	licenses, err := soupnazi.ParseLicenses(lfile, c.stream)
	if err != nil {
		return err
	}

	if c.Raw {
		for _, s := range licenses {
			c.stream.Infof("License: '%s'", s)
			fmt.Printf(" * %s\n", s)
		}
	} else {
		apps := AppCollection{}
		for _, s := range licenses {
			c.stream.Infof("License: '%s'", s)
			token := soupnazi.RawToken(s)
			app, ok := token.Claims["app"].(string)
			if !ok {
				return fmt.Errorf("Couldn't extract application for %v", token.Claims)
			}
			a := apps.App(app)

			feature, ok := token.Claims["f"].(string)
			if !ok {
				return fmt.Errorf("Couldn't extract feature for %v", token.Claims)
			}

			a.Feature(feature)
		}

		for appname, app := range apps {
			fmt.Printf("Application: %s\n", appname)
			for featname, _ := range *app {
				fmt.Printf("  Feature: %s\n", featname)
			}
		}
	}

	return nil
}

func NewListCommand(stream *logrus.Logger, options *Options) *ListCommand {
	return &ListCommand{
		options: options,
		stream:  stream,
	}
}
