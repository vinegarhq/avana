package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	cs "github.com/apprehensions/rbxweb/clientsettings"
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

	var bt cs.BinaryType
	switch os.Args[1] {
	case "player":
		bt = cs.WindowsPlayer
	case "studio":
		bt = cs.WindowsStudio64
	default:
		usage()
	}

	b := binary.New()
	if err := b.Setup(bt); err != nil {
		log.Fatal(err)
	}

	if err := b.Run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}

	slog.Info("Goodbye!")
}
