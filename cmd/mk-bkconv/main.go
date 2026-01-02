package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/galpt/mk-bkconv/pkg/convert"
	"github.com/galpt/mk-bkconv/pkg/kotatsu"
	"github.com/galpt/mk-bkconv/pkg/mihon"
)

func main() {
	// args excludes program name
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	// allow global flags (such as --allow-fallback) to appear anywhere
	allowSourcesFallback := slices.Contains(os.Args, "--allow-fallback")

	// find the subcommand if it's present anywhere among the args
	var sub string
	subIndex := -1
	for i, a := range args {
		if a == "mihon-to-kotatsu" || a == "kotatsu-to-mihon" {
			sub = a
			subIndex = i
			break
		}
	}
	// If no explicit subcommand provided, try to auto-detect from flags (allow running without subcommand)
	if sub == "" {
		// manual scan for -in / --in and global flags when no explicit subcommand provided
		var detIn string
		for i := 0; i < len(args); i++ {
			a := args[i]
			if a == "--allow-fallback" || a == "-allow-fallback" {
				allowSourcesFallback = true
			}
			if a == "-in" || a == "--in" {
				if i+1 < len(args) {
					detIn = args[i+1]
				}
			} else if strings.HasPrefix(a, "-in=") || strings.HasPrefix(a, "--in=") {
				parts := strings.SplitN(a, "=", 2)
				if len(parts) == 2 {
					detIn = parts[1]
				}
			}
		}
		if detIn != "" {
			lower := strings.ToLower(detIn)
			if strings.HasSuffix(lower, ".zip") {
				sub = "kotatsu-to-mihon"
			} else if strings.HasSuffix(lower, ".tachibk") {
				sub = "mihon-to-kotatsu"
			}
		}
		if sub == "" {
			usage()
			os.Exit(1)
		}
	}

	// Build args slice for flag parsing by removing the subcommand token (if present)
	var parsedArgs []string
	if subIndex >= 0 {
		parsedArgs = append([]string{}, args[:subIndex]...)
		parsedArgs = append(parsedArgs, args[subIndex+1:]...)
	} else {
		parsedArgs = append([]string{}, args...)
	}
	// Remove global-only flags (e.g., --allow-fallback or -allow-fallback)
	filteredArgs := make([]string, 0, len(parsedArgs))
	for _, a := range parsedArgs {
		if a == "--allow-fallback" || a == "-allow-fallback" {
			continue
		}
		filteredArgs = append(filteredArgs, a)
	}
	switch sub {
	case "mihon-to-kotatsu":
		fs := flag.NewFlagSet("mihon-to-kotatsu", flag.ExitOnError)
		in := fs.String("in", "", "input mihon backup file (.tachibk)")
		out := fs.String("out", "", "output kotatsu zip file")
		fs.Parse(filteredArgs)
		if *in == "" || *out == "" {
			usage()
			os.Exit(2)
		}
		b, err := mihon.LoadBackup(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading mihon backup: %v\n", err)
			os.Exit(3)
		}
		kb := convert.MihonToKotatsu(b)
		if err := kotatsu.WriteKotatsuZip(*out, kb); err != nil {
			fmt.Fprintf(os.Stderr, "error writing kotatsu zip: %v\n", err)
			os.Exit(4)
		}
		fmt.Println("Conversion complete.")

	case "kotatsu-to-mihon":
		fs := flag.NewFlagSet("kotatsu-to-mihon", flag.ExitOnError)
		in := fs.String("in", "", "input kotatsu zip file")
		out := fs.String("out", "", "output mihon backup file (.tachibk)")
		fs.Parse(filteredArgs)
		if *in == "" || *out == "" {
			usage()
			os.Exit(2)
		}
		kb, err := kotatsu.LoadKotatsuZip(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading kotatsu zip: %v\n", err)
			os.Exit(3)
		}
		b, err := convert.KotatsuToMihon(kb, allowSourcesFallback)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting kotatsu to mihon: %v\n", err)
			os.Exit(5)
		}
		if err := mihon.WriteBackup(*out, b); err != nil {
			fmt.Fprintf(os.Stderr, "error writing mihon backup: %v\n", err)
			os.Exit(4)
		}
		fmt.Println("Conversion complete.")

	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("mk-bkconv: convert between Mihon and Kotatsu backups")
	fmt.Println("USAGE:")
	fmt.Println("  mk-bkconv <mihon-to-kotatsu|kotatsu-to-mihon> -in <input> -out <output> --allow-fallback")
	fmt.Println("    --allow-fallback   this flag allows you to fallback to hashing when there was no mapping for a source found")

}
