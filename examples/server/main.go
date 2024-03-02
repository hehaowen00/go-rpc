package main

import (
	"os"
	"os/signal"

	"github.com/hehaowen00/go-rpc"
	"github.com/hehaowen00/go-rpc/examples/api"
)

func main() {
	notes := NotesService{}

	service := rpc.NewService(api.NotesService, "0.0.0.0:8080", &notes)
	service.Run()
	defer service.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
