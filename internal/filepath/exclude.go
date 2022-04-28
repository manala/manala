package filepath

import "path/filepath"

func Exclude(path string) bool {
	switch filepath.Base(path) {
	case
		// Git
		".git",
		".github",
		// NodeJS
		"node_modules",
		// Composer
		"vendor",
		// IntelliJ
		".idea",
		// Manala
		".manala":
		return true
	}
	return false
}
