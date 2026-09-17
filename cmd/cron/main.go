package main

import (
	"flag"
	"fmt"
	"homelab-sync/internal/backup"
	"os"
)

func main() {
	path := flag.String("path", ".", "path to backup")
	flag.Parse()

	if err := backup.Run(*path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
}
