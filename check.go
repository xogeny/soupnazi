package check

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type LM struct {
	appname string
	secret  string
	jwts    []string
	stream  *logrus.Logger
}

func NewLM(appname string, secret string, level logrus.Level) *LM {
	jwts := []string{}

	// Name of config file (without extension)
	viper.SetConfigName(fmt.Sprintf("%s_licenses", appname))
	viper.SetEnvPrefix(appname)
	viper.AutomaticEnv()

	// Find and read the config file
	err := viper.ReadInConfig()
	if err == nil {
		// Get JWTs
	}

	jwts = viper.GetStringSlice("licenses")

	stream := logrus.New()
	stream.Level = level
	stream.Out = os.Stderr

	stream.Infof("licenses: %v", jwts)

	return &LM{
		appname: appname,
		secret:  secret,
		jwts:    jwts,
		stream:  stream,
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
		token, err := jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			lm.stream.Infof("  Unverified token: %v", token)
			kid, exists := token.Header["kid"]
			if !exists {
				return []byte{}, fmt.Errorf("Verification key not found in token: %v", token)
			}
			key, ok := kid.(string)
			if ok {
				return []byte(key), nil
			} else {
				return []byte{},
					fmt.Errorf("Expected verification key to be a string, but found %T", token)
			}
		})

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
