package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"

	"github.com/xogeny/soupnazi"
)

type Options struct {
	Positional struct {
		Application string `description:"Application Id" required:"true"`
		Feature     string `description:"User Id" required:"true"`
		Secret      string `description:"Secret key" require:"true"`
	} `positional-args:"true"`
	Days      int               `long:"days" short:"d" description:"Days token is valid for"`
	Params    map[string]string `long:"param" description:"String parameters"`
	MAC       string            `long:"mac" short:"m" description:"MAC Address (or 'Any')"`
	Verbose   bool              `long:"verbose" short:"v" description:"Verbose output"`
	Unlimited bool              `long:"unlimited" short:"u" description:"Token doesn't expire"`
	Expired   bool              `long:"expired" short:"e" description:"Token that has already expired"`
}

func main() {
	var Options Options
	parser := flags.NewParser(&Options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		fmt.Printf("Error parsing options: %v\n", err)
		os.Exit(1)
	}

	stream := logrus.New()
	stream.Level = logrus.WarnLevel
	if Options.Verbose {
		stream.Level = logrus.DebugLevel
	}
	stream.Out = os.Stderr

	if Options.Positional.Application == "" {
		fmt.Printf("Error, no application specified\n")
		os.Exit(1)
	}

	if Options.Positional.Feature == "" {
		fmt.Printf("Error, no feature specified\n")
		os.Exit(1)
	}

	if Options.Positional.Secret == "" {
		fmt.Printf("Error, no secret key provided\n")
		os.Exit(1)
	}

	mac := Options.MAC
	if mac == "" {
		fmt.Printf("Error, no MAC address provided\n")
		os.Exit(1)
	}
	if mac == "Any" || mac == "any" {
		mac = "*"
	}

	if Options.Days < 0 {
		fmt.Printf("Error, number of days until expiration must be positive, was %d\n",
			Options.Days)
		os.Exit(1)
	}

	mins := Options.Days * 24 * 60

	// If no days are specified,
	if Options.Days == 0 {
		if Options.Expired {
			mins = 0
		} else if Options.Unlimited {
			mins = -1
		} else {
			fmt.Printf("Either provide number of days to expiration (--days) or (--unlimited)\n")
			os.Exit(1)
		}
	}

	nl := soupnazi.NodeLocked{
		Application: Options.Positional.Application,
		Feature:     Options.Positional.Feature,
		Secret:      Options.Positional.Secret,
		MAC:         mac,
		Minutes:     mins,
		Params:      Options.Params,
	}

	tokenString, err := soupnazi.GenerateNodeLocked(nl, stream)
	if err != nil {
		fmt.Printf("Error generating JWT token: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", tokenString)
	os.Exit(0)
}
