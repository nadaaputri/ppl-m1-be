package utils

import (
	"regexp"
	"strings"
	"time"
)

var (
	CodePattern = regexp.MustCompile(`^[a-z0-9]+([_-][a-z0-9]+)*$`)
	WIB         = time.FixedZone("WIB", 7*60*60)
)

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func EscapeLike(s string) string {
	return likeEscaper.Replace(s)
}
