package main

import (
	cmd "github.com/cacoco/codemetagenerator/cmd"
)

var (
	version = "development"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
