package main

import (
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/jessevdk/go-flags"
)

// Universal flags
type Options struct {
	Verbose bool `long:"verbose" short:"v" description:"Verbose output"`
}

func (o Options) Init(stream *logrus.Logger) error {
	if o.Verbose {
		stream.Level = logrus.DebugLevel
	}

	return nil
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.Default)

	stream := logrus.New()
	stream.Level = logrus.WarnLevel
	stream.Out = os.Stderr

	// Create a new application suite
	parser.AddCommand("add",
		"Add a license",
		"Add a license",
		NewAddCommand(stream, &options))

	parser.AddCommand("list",
		"List available licenses",
		"List available licenses",
		NewListCommand(stream, &options))

	parser.AddCommand("trim",
		"Trim expired licenses",
		"Trim expired licenses",
		NewListCommand(stream, &options))

	parser.AddCommand("rm",
		"Remove a license",
		"Remove a license",
		NewRemoveCommand(stream, &options))

	if _, err := parser.Parse(); err != nil {
		stream.Errorf("Error parsing options: %v", err)
		os.Exit(1)
	}
}
