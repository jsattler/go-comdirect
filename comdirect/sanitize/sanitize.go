package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	digits = regexp.MustCompile(`[0-9]`)
)

func ReplaceDigits(target string, replaceWith string) string {
	return digits.ReplaceAllString(target, replaceWith)
}

func KeepLastN(target string, replaceWith string, n int) string {
	if len(target) < n {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.Repeat(replaceWith, len(target)-n), target[len(target)-n:])
}
