package route

import (
	"regexp"
	"stuber/internal/server/stub"
)

type Route struct {
	Pattern *regexp.Regexp
	Stub    *stub.Stub
}
