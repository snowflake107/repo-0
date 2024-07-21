package strutil

import (
	"github.com/samber/lo"
	"strings"
)

func TrimMultilineWhitespace(input string) string {
	return strings.Join(lo.FilterMap(strings.Split(input, "\n"), func(item string, index int) (string, bool) {
		trimmed := strings.TrimSpace(item)
		return trimmed, trimmed != ""
	}), "\n")
}
