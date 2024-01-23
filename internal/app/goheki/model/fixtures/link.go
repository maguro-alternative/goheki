package fixtures

import (
	"context"
	"testing"
)

type Link struct {
	ID       int64  `db:"id"`
	EntryID  int64  `db:"entry_id"`
	Type     string `db:"type"`
	URL      string `db:"url"`
	Nsfw     bool   `db:"nsfw"`
	Darkness bool   `db:"darkness"`
}

func NewLink(ctx context.Context, setter ...func(l *Link)) *ModelConnector {
	link := &Link{
		Type:     "blog",
		URL:      "https://example.com",
		Nsfw:     false,
		Darkness: false,
	}

	//setter(link)

	return &ModelConnector{
		Model: link,
		setter: func() {
			for _, s := range setter {
				s(link)
			}
		},
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
			// 連番されるIDをセットする
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO link (
					entry_id,
					type,
					url,
					nsfw,
					darkness
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5
				) RETURNING id`,
				link.EntryID,
				link.Type,
				link.URL,
				link.Nsfw,
				link.Darkness,
			).Scan(&link.ID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
