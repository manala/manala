//go:build web_app_build && !web_app_build_embed

package web

import (
	"net/http"
)

func AppFS() (http.FileSystem, string) {
	return http.Dir("web/app/build"), "build"
}
