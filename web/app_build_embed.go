//go:build web_app_build_embed && !web_app_build

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed app/build
var appBuildFS embed.FS

func AppFS() (http.FileSystem, string) {
	appFS, _ := fs.Sub(appBuildFS, "app/build")

	return http.FS(appFS), "build_embed"
}
