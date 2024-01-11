package fixtures

import (
	"context"
	"testing"
)

type Tag struct {
	ID   *int64 `db:"id"`
	Name string `db:"name"`
}

func NewTag(ctx context.Context) *ModelConnector {
	tag := &Tag{
		Name: "tag",
	}

	return &ModelConnector{
		Model: tag,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Tags = append(f.Tags, tag)
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO tag (name) VALUES (:name)", tag)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			tag.ID = &id
		},
	}
}
