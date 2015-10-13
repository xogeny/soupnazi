package soupnazi

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

// LoadJWTs loads all known JWTs
func LoadJWTs(stream *logrus.Logger) []string {
	ret := []string{}

	env := os.Getenv("SOUPNAZI_LICENSES")
	parts := strings.Split(env, ":")

	for _, part := range parts {
		ret = append(ret, part)
	}

	return ret
}
