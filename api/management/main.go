package main

import (
	"log"
	"net/http"
	// nolint:gosec
	_ "net/http/pprof"
	"os"

	"github.com/ethancls/cosmos/management/cmd"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
