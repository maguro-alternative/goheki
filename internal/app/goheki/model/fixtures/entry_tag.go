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

func NewEntryTag(ctx context.Context) *ModelConnector {
	entryTag := &EntryTag{}

	return &ModelConnector{
		Model: entryTag,
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
			result, err := f.dbv1.NamedExecContext(ctx, "INSERT INTO entry_tag (entry_id, tag_id) VALUES (:entry_id, :tag_id)", entryTag)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			entryTag.ID = &id
		},
	}
}
