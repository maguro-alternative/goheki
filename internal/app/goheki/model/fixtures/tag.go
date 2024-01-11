package fixtures

import (
	"context"
	"testing"
)

type Tag struct {
	ID   *int64 `db:"id"`
	Name string `db:"name"`
}

func NewTag(ctx context.Context, setter func(t *Tag)) *ModelConnector {
	tag := &Tag{
		Name: "tag",
	}

	setter(tag)

	return &ModelConnector{
		Model: tag,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Tags = append(f.Tags, tag)
		},
		insertTable: func(t *testing.T, f *Fixture) {
			// 連番されるIDをセットする
			result := f.DBv1.QueryRowxContext(
				ctx,
				"INSERT INTO tag (name) VALUES ($1)",
				tag.Name,
			).Scan(&tag.ID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
