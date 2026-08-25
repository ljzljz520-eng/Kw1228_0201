package main

import (
	"context"
	"flag"
	"log"
	"os"

	"example.com/energycore/internal/flow017"
	"example.com/energycore/internal/httpapi"
)

func main() {
	address := flag.String("addr", ":8080", "HTTP listen address")
	database := flag.String("db", "./energycore.db", "bbolt database path")
	flag.Parse()
	runtime, err := flow017.OpenRuntime(*database)
	if err != nil {
		log.Fatal(err)
	}
	defer runtime.Close()
	server := httpapi.New(runtime)
	ctx, stop := signalContext()
	defer stop()
	log.Printf("energy-core registry listening on %s", *address)
	if err := server.Serve(ctx, *address); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
