package convert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/galpt/mk-bkconv/proto/mihon"
)

// FilterBackupToCommon removes mangas and sources from the Mihon backup
// that don't have matching sources available in both Kotatsu and Mihon.
// It attempts to discover Mihon extension names from a references folder
// (ENV "REFERENCES_ROOT" or ../references by default). If discovery fails
// it falls back to KnownSourceMapping as a conservative whitelist.
func FilterBackupToCommon(b *pb.Backup, kotatsuRawSources []byte) {
	// Discover mihon sources from references (best-effort)
	refRoot := os.Getenv("REFERENCES_ROOT")
	if refRoot == "" {
		// try relative ../references
		cwd, err := os.Getwd()
		if err == nil {
			candidate := filepath.Join(cwd, "..", "..", "references")
			if _, err := os.Stat(candidate); err == nil {
				refRoot = candidate
			}
		}
	}

	mihonNames := make(map[string]struct{})
	// Seed from KnownSourceMapping values (guaranteed known mappings)
	for _, m := range KnownSourceMapping {
		mihonNames[strings.ToLower(m.MihonName)] = struct{}{}
	}

	// If we have a references root, try to walk and discover extension directories
	if refRoot != "" {
		_ = filepath.Walk(refRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			// look for Kotlin/Java extension files path that include the tachiyomi extension package
			p := filepath.ToSlash(path)
			if strings.Contains(p, "/eu/kanade/tachiyomi/extension/") && strings.HasSuffix(p, ".kt") {
				// attempt to extract extension short name (the parent folder under extension)
				parts := strings.Split(p, "/eu/kanade/tachiyomi/extension/")
				if len(parts) > 1 {
					rest := parts[1]
					segs := strings.Split(rest, "/")
					if len(segs) > 0 {
						mihonNames[strings.ToLower(segs[0])] = struct{}{}
					}
				}
			}
			return nil
		})
	}

	// Build allowed ID set from mihonNames using GenerateMihonSourceID where possible
	allowedIDs := make(map[int64]struct{})
	for _, m := range KnownSourceMapping {
		if _, ok := mihonNames[strings.ToLower(m.MihonName)]; ok {
			id := GenerateMihonSourceID(m.MihonName, m.MihonLang, m.MihonVersionID)
			allowedIDs[id] = struct{}{}
		}
	}

	// parse kotatsuRawSources if available to extract kotatsu source names (best-effort)
	kotatsuNames := make(map[string]struct{})
	if len(kotatsuRawSources) > 0 {
		var arr []map[string]interface{}
		if err := json.Unmarshal(kotatsuRawSources, &arr); err == nil {
			for _, el := range arr {
				if s, ok := el["name"].(string); ok {
					kotatsuNames[strings.ToLower(s)] = struct{}{}
				}
			}
		}
	}

	// If allowedIDs is empty, fall back to allowing all KnownSourceMapping IDs
	if len(allowedIDs) == 0 {
		for k := range KnownSourceMapping {
			m := KnownSourceMapping[k]
			id := GenerateMihonSourceID(m.MihonName, m.MihonLang, m.MihonVersionID)
			allowedIDs[id] = struct{}{}
		}
	}

	// Filter BackupManga entries: keep only those whose Source ID is in allowedIDs
	var kept []*pb.BackupManga
	for _, m := range b.BackupManga {
		if _, ok := allowedIDs[m.GetSource()]; ok {
			kept = append(kept, m)
			continue
		}
		// If source not by id, try to match by name from BackupSources
		foundByName := false
		for _, s := range b.BackupSources {
			if s.GetSourceId() == m.GetSource() {
				if s.GetName() != "" {
					if _, ok := mihonNames[strings.ToLower(s.GetName())]; ok {
						foundByName = true
						break
					}
				}
			}
		}
		if foundByName {
			kept = append(kept, m)
		}
	}
	b.BackupManga = kept

	// Filter BackupSources similarly
	var keptSources []*pb.BackupSource
	for _, s := range b.BackupSources {
		if _, ok := allowedIDs[s.GetSourceId()]; ok {
			keptSources = append(keptSources, s)
			continue
		}
		if s.GetName() != "" {
			if _, ok := mihonNames[strings.ToLower(s.GetName())]; ok {
				keptSources = append(keptSources, s)
			}
		}
	}
	b.BackupSources = keptSources
}

// FilterMihonForKotatsu removes Mihon backup entries that don't have a corresponding
// Kotatsu source available. It attempts to discover Kotatsu parser names from
// references (ENV "REFERENCES_ROOT" or ../references by default) and falls back
// to KnownSourceMapping keys if discovery fails.
func FilterMihonForKotatsu(b *pb.Backup) {
	refRoot := os.Getenv("REFERENCES_ROOT")
	if refRoot == "" {
		cwd, err := os.Getwd()
		if err == nil {
			candidate := filepath.Join(cwd, "..", "..", "references")
			if _, err := os.Stat(candidate); err == nil {
				refRoot = candidate
			}
		}
	}

	kotatsuNames := make(map[string]struct{})
	// Seed from KnownSourceMapping keys
	for k := range KnownSourceMapping {
		kotatsuNames[strings.ToLower(k)] = struct{}{}
	}

	if refRoot != "" {
		// look for kotatsu-parsers-master repo
		_ = filepath.Walk(refRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() {
				return nil
			}
			p := filepath.ToSlash(path)
			// match paths like .../kotatsu-parsers-master/src/main/kotlin/.../site/<sourcename>
			if strings.Contains(p, "/kotatsu-parsers-master/src/main/kotlin/") && strings.Contains(p, "/site/") {
				// extract after /site/
				parts := strings.Split(p, "/site/")
				if len(parts) > 1 {
					segs := strings.Split(parts[1], "/")
					if len(segs) > 0 && segs[0] != "" {
						kotatsuNames[strings.ToLower(segs[0])] = struct{}{}
					}
				}
			}
			return nil
		})
	}

	// Build allowed Mihon IDs for kotatsu-supported sources via KnownSourceMapping
	allowedIDs := make(map[int64]struct{})
	for k := range KnownSourceMapping {
		if _, ok := kotatsuNames[strings.ToLower(k)]; ok {
			m := KnownSourceMapping[k]
			id := GenerateMihonSourceID(m.MihonName, m.MihonLang, m.MihonVersionID)
			allowedIDs[id] = struct{}{}
		}
	}

	// If allowedIDs empty, keep existing backup untouched (conservative)
	if len(allowedIDs) == 0 {
		return
	}

	// Filter BackupManga and BackupSources
	var kept []*pb.BackupManga
	for _, m := range b.BackupManga {
		if _, ok := allowedIDs[m.GetSource()]; ok {
			kept = append(kept, m)
		}
	}
	b.BackupManga = kept

	var keptSources []*pb.BackupSource
	for _, s := range b.BackupSources {
		if _, ok := allowedIDs[s.GetSourceId()]; ok {
			keptSources = append(keptSources, s)
		}
	}
	b.BackupSources = keptSources
}
