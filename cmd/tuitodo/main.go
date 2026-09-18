package main

import (
	"fmt"
	"os"

	"tuitodo/internal/store"
	"tuitodo/internal/ui"
)

func main() {
	databasePath, err := store.DefaultDatabasePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "database path: %v\n", err)
		os.Exit(1)
	}
	sqliteStore, err := store.Open(databasePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer sqliteStore.Close()

	if err := ui.Run(sqliteStore, sqliteStore); err != nil {
		fmt.Fprintf(os.Stderr, "run ui: %v\n", err)
		os.Exit(1)
	}
}
