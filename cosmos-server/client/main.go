package main

import (
	"os"

	"github.com/ethancls/cosmos-server/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
