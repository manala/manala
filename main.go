package main

import (
	"manala/cmd"
	"os"
)

// Set at build time, by goreleaser, via ldflags
var version = "dev"

func main() {
	cmd.Execute(version, os.Stdin, os.Stdout, os.Stderr)
}
