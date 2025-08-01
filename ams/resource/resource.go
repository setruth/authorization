package resource

import (
	"embed"
	"io/fs"
)

//go:embed web
var embedFS embed.FS
var WebStatic fs.FS

func init() {
	var err error
	WebStatic, err = fs.Sub(embedFS, "web")
	if err != nil {
		panic(err)
	}
}
