package main

import (
	"manala/cmd"
	"os"
)

// Set at build time, by goreleaser, via ldflags
var version = "dev"

// Default repository
var defaultRepository = "https://github.com/manala/manala-recipes.git"

func main() {
	cmd.Execute(version, defaultRepository, os.Stdout, os.Stderr)
}
