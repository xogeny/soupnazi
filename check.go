package soupnazi

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

type LM struct {
	appname string
	secret  string
	jwts    []string
	stream  *logrus.Logger
}

// NewLM creates a new license manager instance.  The first argument is the name
// of the application that is running.  The second is a shared secret (shared between
// the running application and whoever is generating licenses for the app).  Finally,
// an optional loglevel is included.  Technically, the log level is zero or more
// log levels but only the last one is used.
func NewLM(appname string, secret string, levels ...logrus.Level) *LM {
	stream := logrus.New()
	stream.Level = logrus.WarnLevel
	for _, level := range levels {
		stream.Level = level
	}
	stream.Out = os.Stderr

	jwts := LoadJWTs(stream)

	stream.Infof("licenses: %v", jwts)

	return &LM{
		appname: appname,
		secret:  secret,
		jwts:    jwts,
		stream:  stream,
	}
}

func RawToken(j string) *jwt.Token {
	var ret *jwt.Token = nil

	// The last part of the key should be the base 64 encoding of a sha256 hash.
	// If it isn't 43 characters, then something is wrong.
	parts := strings.Split(j, ".")
	if len(parts[2]) != 43 {
		return nil
	}

	// TODO: This needs to do a more thorough integrity check.  It
	jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
		// If we got here, the basic checks passed
		ret = token

		return nil, nil
	})

	return ret
}

func SyntaxCheck(j string) bool {
	return RawToken(j) != nil
}

func KeyFunc(key string, stream *logrus.Logger) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		alg, exists := token.Header["alg"]
		if !exists {
			return nil, fmt.Errorf("Algorithm missing from header")
		}
		algs, ok := alg.(string)
		if !ok {
			return nil, fmt.Errorf("Algorithm not a string")
		}
		stream.Infof("  Algorithm: %s", algs)
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if algs != "HS256" {
			return nil, fmt.Errorf("Unexpected algorithm: %s", algs)
		}

		stream.Infof("  Unverified token: %v", token)

		return []byte(KeyHash(key)), nil
	}
}

// License checks to see if a given feature is licensed and, if it is,
// it provides a map of parameters about the feature.  These parameters
// generally represent any limitations that might be associated with the
// feature.  The parameters are feature specific and are encoded in
// the JWTs themselves
func (lm LM) License(feature string) (map[string]string, error) {
	blank := map[string]string{}

	// Loop over JWTs
	for _, j := range lm.jwts {
		lm.stream.Debugf("Processing JWT: '%s'", j)
		// Parse the JWTs
		token, err := jwt.Parse(j, KeyFunc(lm.secret, lm.stream))

		// If there was an error parsing the token, skip it
		if err != nil {
			lm.stream.Debugf("  Error parsing/validating: %v", err)
			continue
		}

		if !token.Valid {
			lm.stream.Debugf("  Invalid token")
			continue
		}

		// TODO: Check if it unlocks feature
		if app, ok := token.Claims["app"]; ok {
			if app != lm.appname {
				lm.stream.Debugf("  Token is for application %s, not %s", app, lm.appname)
				continue
			}
			if feat, ok := token.Claims["f"]; ok {
				if feat != feature {
					lm.stream.Debugf("  Token is for feature %s, not %s", feat, feature)
					continue
				}
				lm.stream.Debugf("  Token matches")

				raw, ok := token.Claims["p"].(map[string]interface{})
				params := map[string]string{}
				if ok {
					lm.stream.Infof("  Params: %v (%T)", raw, raw)
					for k, v := range raw {
						lm.stream.Infof("    Parameter %s = %v", k, v)
						params[k] = fmt.Sprintf("%v", v)
					}
				} else {
					lm.stream.Warnf("  Token included invalid parameter field")
				}

				// License found
				lm.stream.Infof("  Parameter set: %v", params)
				return params, nil
			} else {
				lm.stream.Debugf("  Token does not specify a feature")
			}
		} else {
			lm.stream.Debugf("  Token does not specify an applications")
		}
	}

	// No soup for you!
	return blank, fmt.Errorf("Unable to locate a license for feature %s of application %s",
		feature, lm.appname)
}
