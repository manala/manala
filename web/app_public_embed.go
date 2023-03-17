//go:build !web_app_build && !web_app_build_embed

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed app/public
var appPublicFS embed.FS

func AppFS() (http.FileSystem, string) {
	appFS, _ := fs.Sub(appPublicFS, "app/public")

	return http.FS(appFS), "public_embed"
}
