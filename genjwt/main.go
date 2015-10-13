package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/dgrijalva/jwt-go"
)

type Options struct {
	Positional struct {
		Application string `description:"Application Id" required:"true"`
		Feature     string `description:"User Id" required:"true"`
		Secret      string `description:"Secret key" require:"true"`
		Params      string `description:"Parameters (as JSON)" require:"false"`
	} `positional-args:"true"`
}

func main() {
	var Options Options
	parser := flags.NewParser(&Options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		log.Printf("Error parsing options: %v", err)
		os.Exit(1)
	}

	if Options.Positional.Application == "" {
		log.Printf("Error, no application specified")
		os.Exit(1)
	}

	if Options.Positional.Feature == "" {
		log.Printf("Error, no feature specified")
		os.Exit(1)
	}

	if Options.Positional.Secret == "" {
		log.Printf("Error, no secret key provided")
		os.Exit(1)
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	var key = Options.Positional.Secret

	token.Header["kid"] = key

	// Set some claims
	token.Claims["app"] = Options.Positional.Application
	token.Claims["f"] = Options.Positional.Feature
	// TODO: Make this optional with a duration flag
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	if Options.Positional.Params != "" {
		params := map[string]interface{}{}
		err := json.Unmarshal([]byte(Options.Positional.Params), &params)
		if err == nil {
			token.Claims["p"] = params
		} else {
			log.Printf("Error parsing parameters: %v", err)
			os.Exit(1)
		}
	}

	log.Printf("Claims: %v", token.Claims)
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", tokenString)
	os.Exit(0)
}
