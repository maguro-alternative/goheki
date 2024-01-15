package fixtures

import (
	//"github.com/maguro-alternative/goheki/pkg/db"

	"context"
	"testing"
)

type Source struct {
	ID   *int64 `db:"id"`
	Name string `db:"name"`
	Url  string `db:"url"`
	Type string `db:"type"`
}

func NewSource(ctx context.Context, setter ...func(s *Source)) *ModelConnector {
	source := &Source{
		Name: "test",
		Url:  "https://example.com",
		Type: "2",
	}

	//setter(source)

	return &ModelConnector{
		Model: source,
		setter: func() {
			for _, s := range setter {
				s(source)
			}
		},
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
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO source (
					name,
					url,
					type
				) VALUES (
					$1,
					$2,
					$3
				) RETURNING id`,
				source.Name,
				source.Url,
				source.Type,
			).Scan(&source.ID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
			// 連番されるIDをセットする
		},
	}
}
