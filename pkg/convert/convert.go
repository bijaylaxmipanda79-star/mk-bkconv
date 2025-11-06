package convert

import (
	"github.com/galpt/mk-bkconv/pkg/kotatsu"
	pb "github.com/galpt/mk-bkconv/proto/mihon"
)

// Helper functions to work with optional string pointers
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// MihonToKotatsu converts from protobuf-based Mihon backup to Kotatsu backup
func MihonToKotatsu(b *pb.Backup) *kotatsu.KotatsuBackup {
	kb := &kotatsu.KotatsuBackup{}

	for i, m := range b.BackupManga {
		fav := kotatsu.KotatsuFavouriteEntry{
			MangaId:    int64(i + 1),
			CategoryId: 0, // Will be updated if manga has categories
			SortKey:    i,
			Pinned:     false,
			CreatedAt:  m.DateAdded,
			Manga: kotatsu.KotatsuManga{
				Id:         int64(i + 1),
				Title:      m.Title,
				Url:        m.Url,
				PublicUrl:  m.Url,
				CoverUrl:   stringVal(m.ThumbnailUrl),
				LargeCover: stringVal(m.ThumbnailUrl),
				Author:     stringVal(m.Author),
				Source:     "",
				Tags:       []interface{}{},
			},
		}

		// Assign first category if exists
		if len(m.Categories) > 0 {
			fav.CategoryId = m.Categories[0]
		}

		kb.Favourites = append(kb.Favourites, fav)
	}

	// Convert categories
	for _, c := range b.BackupCategories {
		kb.Categories = append(kb.Categories, kotatsu.KotatsuCategory{
			CategoryId: c.Id,
			CreatedAt:  c.Order,
			SortKey:    0,
			Title:      c.Name,
		})
	}

	return kb
}

// KotatsuToMihon converts from Kotatsu backup to protobuf-based Mihon backup
func KotatsuToMihon(kb *kotatsu.KotatsuBackup) *pb.Backup {
	b := &pb.Backup{}

	// Build a map of manga ID -> chapters from the index
	chaptersByManga := make(map[int64][]*pb.BackupChapter)
	for _, idx := range kb.Index {
		var chapters []*pb.BackupChapter
		for _, kc := range idx.Chapters {
			chapters = append(chapters, &pb.BackupChapter{
				Url:           kc.Url,
				Name:          kc.Name,
				Scanlator:     stringPtr(kc.Scanlator),
				Read:          false,
				Bookmark:      false,
				LastPageRead:  0,
				ChapterNumber: kc.Number,
			})
		}
		chaptersByManga[idx.MangaId] = chapters
	}

	// Convert favourites to mangas with their chapters
	for _, fav := range kb.Favourites {
		km := fav.Manga
		m := &pb.BackupManga{
			Source:       0,
			Url:          km.Url,
			Title:        km.Title,
			Author:       stringPtr(km.Author),
			Artist:       stringPtr(""),
			Description:  stringPtr(""),
			Genre:        []string{},
			Status:       0,
			ThumbnailUrl: stringPtr(km.CoverUrl),
			DateAdded:    fav.CreatedAt,
			Chapters:     chaptersByManga[km.Id], // attach chapters from index
			Categories:   []int64{fav.CategoryId},
			Favorite:     true,
		}
		b.BackupManga = append(b.BackupManga, m)
	}

	// Convert categories
	for _, c := range kb.Categories {
		b.BackupCategories = append(b.BackupCategories, &pb.BackupCategory{
			Name:  c.Title,
			Order: c.CreatedAt,
			Id:    c.CategoryId,
		})
	}

	return b
}
