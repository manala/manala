//go:build web_app_build && !web_app_build_embed

package web

import (
	"net/http"
)

func init() {
	appFS.Name = "build"
	appFS.FS = http.Dir("web/app/build")
}
