package fixtures

import (
	"context"
	"testing"
)

type EntryTag struct {
	ID      *int64 `db:"id"`
	EntryID *int64 `db:"entry_id"`
	TagID   *int64 `db:"tag_id"`
}

func NewEntryTag(ctx context.Context, setter ...func(e *EntryTag)) *ModelConnector {
	entryTag := &EntryTag{}

	//setter(entryTag)

	return &ModelConnector{
		Model: entryTag,
		setter: func() {
			for _, s := range setter {
				s(entryTag)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.EntryTags = append(f.EntryTags, entryTag)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				entryTag.EntryID = entry.ID
			case *Tag:
				tag := connectingModel.(*Tag)
				entryTag.TagID = tag.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, entryTag)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO entry_tag (
					entry_id,
					tag_id
				) VALUES (
					$1,
					$2
				) RETURNING id`,
				entryTag.EntryID,
				entryTag.TagID,
			).Scan(&entryTag.ID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
