package fixtures

import (
	"context"
	"testing"
)

type Link struct {
	ID       *int64 `db:"id"`
	EntryID  *int64 `db:"entry_id"`
	Type     string `db:"type"`
	URL      string `db:"url"`
	Nsfw     bool   `db:"nsfw"`
	Darkness bool   `db:"darkness"`
}

func NewLink(ctx context.Context) *ModelConnector {
	link := &Link{
		Type:     "blog",
		URL:      "https://example.com",
		Nsfw:     false,
		Darkness: false,
	}

	return &ModelConnector{
		Model: link,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Links = append(f.Links, link)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				link.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, link)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO link (entry_id, type, url, nsfw, darkness) VALUES (:entry_id, :type, :url, :nsfw, :darkness)", link)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			link.ID = &id
		},
	}
}
