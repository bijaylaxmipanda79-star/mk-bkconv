package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/galpt/mk-bkconv/pkg/mihon"
)

func main() {
	inFile := flag.String("in", "", "input mihon backup file")
	flag.Parse()

	if *inFile == "" {
		log.Fatal("-in required")
	}

	backup, err := mihon.LoadBackup(*inFile)
	if err != nil {
		log.Fatalf("Error loading backup: %v", err)
	}

	fmt.Println("=== BACKUP SOURCES ===")
	fmt.Printf("Total sources: %d\n\n", len(backup.BackupSources))

	for _, src := range backup.BackupSources {
		sid := src.GetSourceId()
		fmt.Printf("Source ID: %d (0x%016x)\n", sid, uint64(sid))
		fmt.Printf("Name: %s\n\n", src.GetName())
	}

	// Show unique sources from manga
	sourceCounts := make(map[int64]int)
	for _, m := range backup.BackupManga {
		sourceCounts[m.GetSource()]++
	}

	fmt.Println("=== MANGA BY SOURCE ===")
	for srcID, count := range sourceCounts {
		fmt.Printf("Source ID %d: %d manga\n", srcID, count)
	}
}
