package fixtures

import (
	//"github.com/maguro-alternative/goheki/pkg/db"

	"context"
	"testing"
)

type Source struct {
	ID      *int64 `db:"id"`
	Name    string `db:"name"`
	Url     string `db:"url"`
	Type    string `db:"type"`
}

func NewSource(ctx context.Context,) *ModelConnector {
	source := &Source{
		Name: "test",
		Url:  "https://example.com",
		Type: "2",
	}

	return &ModelConnector{
		Model: source,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Sources = append(f.Sources, source)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				entry.SourceID = *source.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, source)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.dbv1.NamedExecContext(ctx, "INSERT INTO source (name, url, type) VALUES (:name, :url, :type)", source)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			source.ID = &id
		},
	}
}