package utils

import (
	"regexp"
	"time"
)

var (
	CodePattern = regexp.MustCompile(`^[a-z0-9]+([_-][a-z0-9]+)*$`)
	WIB         = time.FixedZone("WIB", 7*60*60)
)
