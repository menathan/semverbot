package main

import (
	"github.com/menathan/semverbot/pkg/cli/exec"
)

// main bootstraps the `sbot` CLI app.
func main() {
	_ = exec.Run()
}
