package main

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/internal/server"
)

var errExitCode = 1

func main() {
	err := server.Run()
	if err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(errExitCode)
	}
}