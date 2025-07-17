// This file is kept for backward compatibility
// The actual CLI application is in cmd/cli/main.go
// To build the CLI: go build ./cmd/cli
// To build the webhook server: go build ./cmd/webhook-server

package main

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/internal/cmd"
)

var errExitCode = 1

func main() {
	fmt.Println("Warning: This entry point is deprecated. Please use ./cmd/cli instead.")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(errExitCode)
	}
}
