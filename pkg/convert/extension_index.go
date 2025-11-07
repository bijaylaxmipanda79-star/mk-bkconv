package convert

// ExtensionMetadata represents information about a Mihon extension
type ExtensionMetadata struct {
	PackageName string // e.g., "eu.kanade.tachiyomi.extension.en.mangadex"
	Name        string // e.g., "Tachiyomi: MangaDex"
	Lang        string // e.g., "en"
	Version     string // e.g., "1.4.123"
	Sources     []SourceInExtension
}

// SourceInExtension represents a source provided by an extension
type SourceInExtension struct {
	Name    string // e.g., "MangaDex"
	Lang    string // e.g., "en"
	ID      int64  // e.g., 2499283573021220255
	BaseURL string // e.g., "https://mangadex.org"
}

// KeiyoushiIndex caches the Keiyoushi extension index
// In a production tool, this would be fetched from:
// https://raw.githubusercontent.com/keiyoushi/extensions/repo/index.min.json
// and cached locally or embedded in the binary
var KeiyoushiIndex = map[int64]ExtensionMetadata{
	// This would be populated from the index.min.json
	// For now, we'll use the known mappings approach
	// TODO: Implement index fetching and parsing
}

// GetExtensionForSource returns the extension package name for a given source ID
func GetExtensionForSource(sourceID int64) (packageName string, found bool) {
	ext, found := KeiyoushiIndex[sourceID]
	if !found {
		return "", false
	}
	return ext.PackageName, true
}
