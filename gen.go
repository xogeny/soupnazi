package soupnazi

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

type NodeLocked struct {
	Application string
	Feature     string
	Secret      string
	MAC         string
	Minutes     int
	Params      map[string]string
}

func GenerateNodeLocked(details NodeLocked, stream *logrus.Logger) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	var key = details.Secret

	// Set some claims
	token.Claims["app"] = details.Application
	token.Claims["f"] = details.Feature

	mins := details.Minutes

	stream.Infof("Generating token:")

	// If a negative number was explicitly provided, no expiration is
	// specified.  Otherwise, the details.Days parameter is the number
	// of days until expiration
	if mins < 0 {
		stream.Infof("  Token will never expire")
	} else {
		dur := time.Minute * time.Duration(mins)
		token.Claims["exp"] = time.Now().Add(dur).Unix()
		stream.Infof("  Token will expire in: %s", dur.String())
	}

	if details.MAC == "*" {
		stream.Infof("  Token will work for: any MAC address")
	} else {
		stream.Infof("  Token will only work for MAC address: %s", details.MAC)
		token.Claims["a"] = details.MAC
	}

	pstr := ""
	sep := ""
	for k, v := range details.Params {
		pstr = pstr + sep + k + "=" + v
		sep = ", "
	}

	stream.Infof("  Params: %s", pstr)
	if len(details.Params) > 0 {
		token.Claims["p"] = details.Params
	}

	stream.Infof("  Claims: %v", token.Claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(key))

	if err == nil {
		stream.Infof("  Token: %s", tokenString)
	} else {
		stream.Infof("  Error creating token: %v", err)
	}

	return tokenString, err
}
