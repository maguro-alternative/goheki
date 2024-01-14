package fixtures

import (
	//"github.com/maguro-alternative/goheki/internal/app/goheki/model/fixtures"

	"context"
	"testing"
	"time"
)

type Entry struct {
	ID        *int64    `db:"id"`
	SourceID  int64     `db:"source_id"`
	Name      string    `db:"name"`
	Image     string    `db:"image"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func NewEntry(ctx context.Context, setter func(e *Entry)) *ModelConnector {
	entry := &Entry{
		SourceID:  1,
		Name:      "test",
		Image:     "https://example.com",
		Content:   "test",
		CreatedAt: time.Now(),
	}

	setter(entry)

	return &ModelConnector{
		Model: entry,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Entrys = append(f.Entrys, entry)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *EntryTag:
				entryTag := connectingModel.(*EntryTag)
				entryTag.EntryID = entry.ID
			case *BWH:
				bwh := connectingModel.(*BWH)
				bwh.EntryID = entry.ID
			case *Personality:
				personality := connectingModel.(*Personality)
				personality.EntryID = entry.ID
			case *Link:
				link := connectingModel.(*Link)
				link.EntryID = entry.ID
			case *HairColor:
				hairColor := connectingModel.(*HairColor)
				hairColor.EntryID = *entry.ID
			case *HairLength:
				hairLength := connectingModel.(*HairLength)
				hairLength.EntryID = entry.ID
			case *HairStyle:
				hairStyle := connectingModel.(*HairStyle)
				hairStyle.EntryID = entry.ID
			case *HekiRadarChart:
				hekiRadarChart := connectingModel.(*HekiRadarChart)
				hekiRadarChart.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, entry)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO entry (
					source_id,
					name,
					image,
					content,
					created_at
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5
				) RETURNING id`,
				entry.SourceID,
				entry.Name,
				entry.Image,
				entry.Content,
				entry.CreatedAt,
			).Scan(&entry.ID)
			// 連番されるIDをセットする
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
