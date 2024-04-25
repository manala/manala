package main

import (
	"manala/app/config"
	"manala/cmd"
	"os"
)

// Set at build time, by goreleaser, via ldflags
var version = "dev"

// Default repository
const repository = "https://github.com/manala/manala-recipes.git"

func main() {
	conf := &config.Config{
		Debug:      false,
		Repository: repository,
		CacheDir:   "",
	}

	code := cmd.Execute(version, conf, os.Stdin, os.Stdout, os.Stderr)

	os.Exit(code)
}
