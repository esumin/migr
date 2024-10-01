package matcher_v1

import (
	"strings"
)

func MatchImport(line string) string {
	if strings.Contains(line, "github.com/pkg/errors") {
		return strings.ReplaceAll(line, "github.com/pkg/errors", "github.com/kanisterio/errkit")
	}
	return ""
}
