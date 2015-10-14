package soupnazi

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"

	jwtgo "github.com/dgrijalva/jwt-go"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestLifecycles(t *testing.T) {
	stream := logrus.New()
	stream.Level = logrus.WarnLevel
	stream.Out = os.Stderr

	exp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiUjJEMiIsImFwcCI6ImFwcGxpY2F0aW9uIiwiZiI6ImZlYXR1cmUifQ.p6h1xDftlnsMtrAbH5VgA9OTq0qAdqEb0YECVqd8ZGg"

	details := NodeLocked{
		Application: "application",
		Feature:     "feature",
		Secret:      "Don't tell anybody",
		MAC:         "R2D2",
		Minutes:     -1,
		Params:      map[string]string{},
	}

	Convey("Should be able to generate a key", t, func(c C) {
		jwt, err := GenerateNodeLocked(details, stream)
		NoError(c, err)
		Equals(c, jwt, exp)
	})

	Convey("Should be able to check (not verify) a valid key", t, func(c C) {
		token := RawToken(exp)
		NotNil(c, token)
	})

	Convey("Should be able to detect missing characters on the front", t, func(c C) {
		token := RawToken(exp[1:])
		Equals(c, token, nil)

		token = RawToken(exp[2:])
		Equals(c, token, nil)
	})

	Convey("Should be able to detect missing characters at the end", t, func(c C) {
		token := RawToken(exp[:(len(exp) - 1)])
		Equals(c, token, nil)

		token = RawToken(exp[:(len(exp) - 2)])
		Equals(c, token, nil)
	})

	Convey("Should be able to generate and then verify a key", t, func(c C) {
		jwt, err := GenerateNodeLocked(details, stream)
		NoError(c, err)

		_, err = jwtgo.Parse(jwt, KeyFunc(details.Secret, stream))
		NoError(c, err)
		Equals(c, jwt, exp)
	})
}