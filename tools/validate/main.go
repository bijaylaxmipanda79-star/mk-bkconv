package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/galpt/mk-bkconv/pkg/mihon"
)

func main() {
	in := flag.String("in", "", "input mihon backup file (.tachibk)")
	flag.Parse()
	if *in == "" {
		log.Fatal("-in required")
	}

	// Use the new protobuf-based loader
	backup, err := mihon.LoadBackup(*in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading backup: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("Validated backup: parse OK\n")
	fmt.Printf("  Manga count: %d\n", len(backup.BackupManga))
	fmt.Printf("  Category count: %d\n", len(backup.BackupCategories))
}
