package fixtures

import (
	//"github.com/maguro-alternative/goheki/internal/app/goheki/model/fixtures"

	"testing"
	"context"
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

func NewEntry(ctx context.Context) *ModelConnector {
	entry := &Entry{
		SourceID: 1,
		Name: "test",
		Image: "https://example.com",
		Content: "test",
		CreatedAt: time.Now(),
	}

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
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, entry)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.dbv1.NamedExecContext(ctx, "INSERT INTO entry (source_id, name, image, content, created_at) VALUES (:source_id, :name, :image, :content, :created_at)", entry)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			entry.ID = &id
		},
	}
}
