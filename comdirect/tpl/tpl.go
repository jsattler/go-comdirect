package tpl

import (
	"embed"
)

//go:embed *.tmpl
var Default embed.FS
