package main

import (
	"github.com/ethancls/cosmos-server/signal/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
