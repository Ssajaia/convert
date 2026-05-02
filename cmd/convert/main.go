package main

import (
	"github.com/convert/internal/cli"
)

// version is injected at build time via -ldflags.
var version = "dev"

func main() {
	cli.Execute(version)
}
