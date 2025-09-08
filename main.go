package main

import (
	"os"

	"github.com/KatrinSalt/backend-challenge-go/cmd/cli"
)

func main() {
	os.Exit(cli.CLI(os.Args))
}
