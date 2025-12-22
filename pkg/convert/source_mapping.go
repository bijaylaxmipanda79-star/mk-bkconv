package convert

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// KnownSourceMapping provides approximate source ID mappings between Kotatsu and Mihon
// for common sources that exist in both ecosystems.
//
// Note: These are APPROXIMATE mappings. Kotatsu and Mihon use completely different
// source implementations and ecosystems. This mapping is for sources that:
// 1. Exist in both ecosystems
// 2. Target the same website
// 3. Have similar enough behavior that a migration makes sense
//
// Users will still need to verify and possibly manually adjust sources after import.
var KnownSourceMapping = map[string]SourceMapping{
	"MANGADEX": {
		MihonName:      "MangaDex",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "Official MangaDex source",
	},
	"MANGAPARK": {
		MihonName:      "MangaPark",
		MihonLang:      "en",
		MihonVersionID: 1,
		Notes:          "MangaPark English",
	},
	// MangaFire doesn't have an official Mihon extension
	// Users will need to find an alternative or use a web source
	"MANGAFIRE_EN": {
		MihonName:      "mangafire",
		MihonLang:      "en",
		MihonVersionID: 1,
		Notes:          "Approximate - verify after import",
	},
	// Add more known mappings here as discovered
	// Bulk entries inferred from local `mihon-extensions-source-main/lib-multisrc`.
	// These entries use the library directory name as the Mihon extension name.
	// The keys are uppercase candidate Kotatsu identifiers; they should be
	// adjusted if your Kotatsu uses different naming.
	"MADARA": {
		MihonName:      "madara",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGABOX": {
		MihonName:      "mangabox",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGACATALOG": {
		MihonName:      "mangacatalog",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGADVENTURE": {
		MihonName:      "mangadventure",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGAESP": {
		MihonName:      "mangaesp",
		MihonLang:      "es",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGAHUB": {
		MihonName:      "mangahub",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGARAW": {
		MihonName:      "mangaraw",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGAREADER": {
		MihonName:      "mangareader",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGATHEMESIA": {
		MihonName:      "mangathemesia",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANGAWORLD": {
		MihonName:      "mangaworld",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MANHWAZ": {
		MihonName:      "manhwaz",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MMLOOK": {
		MihonName:      "mmlook",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MMRCMS": {
		MihonName:      "mmrcms",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MONOCHROME": {
		MihonName:      "monochrome",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MULTICHAN": {
		MihonName:      "multichan",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"PAPRIKA": {
		MihonName:      "paprika",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"PEACHSCAN": {
		MihonName:      "peachscan",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"PIZZAREADER": {
		MihonName:      "pizzareader",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"READERFRONT": {
		MihonName:      "readerfront",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"SCANR": {
		MihonName:      "scanr",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"SENKURO": {
		MihonName:      "senkuro",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"TERRASCAN": {
		MihonName:      "terrascan",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"VERCOMICS": {
		MihonName:      "vercomics",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"WPcomics": {
		MihonName:      "wpcomics",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"YUYU": {
		MihonName:      "yuyu",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"ZBULU": {
		MihonName:      "zbulu",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"ZEISTMANGA": {
		MihonName:      "zeistmanga",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"ZMANGA": {
		MihonName:      "zmanga",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"COMICIVIEWER": {
		MihonName:      "comiciviewer",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"FANSUBSCAT": {
		MihonName:      "fansubscat",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"GIGAVIEWER": {
		MihonName:      "gigaviewer",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"GMANGA": {
		MihonName:      "gmanga",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"MADTHEME": {
		MihonName:      "madtheme",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"LILIANA": {
		MihonName:      "liliana",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
	"HENTAIFOX": {
		MihonName:      "HentaiFox",
		MihonLang:      "all",
		MihonVersionID: 1,
		Notes:          "",
	},
}

// SourceMapping represents a known mapping from Kotatsu to Mihon source
type SourceMapping struct {
	MihonName      string // Exact source name as it appears in Mihon
	MihonLang      string // Language code (e.g., "en", "all")
	MihonVersionID int    // Version ID (usually 1)
	Notes          string // Additional notes for users
}

// GenerateMihonSourceID generates a source ID using Mihon's algorithm:
// MD5("sourcename/lang/versionid")[0:8] as Long with sign bit cleared
func GenerateMihonSourceID(name, lang string, versionID int) int64 {
	// Mihon uses lowercase for the source name
	key := fmt.Sprintf("%s/%s/%d", strings.ToLower(name), lang, versionID)

	// MD5 hash
	hash := md5.Sum([]byte(key))

	// Take first 8 bytes and combine into a Long (same as Mihon's implementation)
	var id int64
	for i := 0; i < 8; i++ {
		id |= int64(hash[i]) << (8 * (7 - i))
	}

	// Clear the sign bit (set MSB to 0) to ensure positive ID
	id &= 0x7FFFFFFFFFFFFFFF

	return id
}

// LookupKnownSource attempts to find a known Mihon mapping for a Kotatsu source
func LookupKnownSource(kotatsuSource string) (sourceID int64, sourceName string, found bool) {
	if mapping, exists := KnownSourceMapping[kotatsuSource]; exists {
		id := GenerateMihonSourceID(mapping.MihonName, mapping.MihonLang, mapping.MihonVersionID)
		return id, mapping.MihonName, true
	}
	return 0, "", false
}
