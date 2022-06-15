package main

import (
	"os"

	"github.com/jonasbleyl/ansible-find/internal/cli"
)

func main() {
	err := cli.Setup().Execute()
	if err != nil {
		os.Exit(1)
	}
}
