package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sewnie/rbxbin"
	"github.com/sewnie/rbxweb"
	"github.com/vinegarhq/avana/internal/binary"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s player|studio [ARGS...]", os.Args[0])
	os.Exit(1)
}

func main() {
	// name-2006-01-02T15:04:05Z07:00.log
	if len(os.Args) < 2 {
		usage()
	}

	var binaryType rbxweb.BinaryType
	switch os.Args[1] {
	case "player":
		binaryType = rbxbin.WindowsPlayer
	case "studio":
		binaryType = rbxbin.WindowsStudio
	default:
		usage()
	}

	client := rbxweb.NewClient()
	deployment, err := rbxbin.GetDeployment(client, binaryType, "")
	if err != nil {
		log.Fatalf("deployment: %s", err)
	}

	b := binary.New(client, deployment)
	
	if err := b.Setup(); err != nil {
		log.Fatalf("setup %s: %s", deployment.GUID, err)
	}

	if err := b.Run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}

	slog.Info("Goodbye!")
}
