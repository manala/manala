//go:build web_app_build_embed && !web_app_build

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed app/build
var appBuildFS embed.FS

func init() {
	fS, _ := fs.Sub(appBuildFS, "app/build")

	appFS.Name = "build_embed"
	appFS.FS = http.FS(fS)
}
