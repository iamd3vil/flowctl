package main

import (
	"embed"

	"github.com/cvhariharan/flowctl/cmd"
)

//go:embed site/build/*
//go:embed site/build/_app
//go:embed configs
//go:embed migrations
var staticFiles embed.FS

func main() {
	cmd.StaticFiles = staticFiles
	cmd.Execute()
}
