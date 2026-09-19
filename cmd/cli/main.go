package main

import (
	"context"
	"fmt"
	"os"

	"homelab-backuper/internal/backup"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <folder> \n", os.Args[0])
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to load .env:", err)
		os.Exit(1)
	}

	folder := os.Args[1]
	bucket := os.Getenv("BACKUP_BUCKET")
	gpgRecipient := os.Getenv("GPG_RECIPIENT")

	ctx := context.Background()

	if err := backup.Run(ctx, folder, bucket, gpgRecipient); err != nil {
		fmt.Fprintln(os.Stderr, "backup failed:", err)
		os.Exit(1)
	}

	fmt.Println("backup complete")
}
