package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jessevdk/go-flags"

	"github.com/dgrijalva/jwt-go"
)

type Options struct {
	Positional struct {
		Application string `description:"Application Id" required:"true"`
		Feature     string `description:"User Id" required:"true"`
		Secret      string `description:"Secret key" require:"true"`
	} `positional-args:"true"`
	Days    int               `long:"days" description:"Number of days token is valid for"`
	Params  map[string]string `long:"param" description:"String parameters"`
	Verbose bool              `long:"verbose" short:"v" description:"Verbose output"`
}

func main() {
	var Options Options
	parser := flags.NewParser(&Options, flags.Default)

	stream := logrus.New()
	stream.Level = logrus.WarnLevel
	if Options.Verbose {
		stream.Level = logrus.DebugLevel
	}
	stream.Out = os.Stderr

	if _, err := parser.Parse(); err != nil {
		stream.Errorf("Error parsing options: %v", err)
		os.Exit(1)
	}
}
