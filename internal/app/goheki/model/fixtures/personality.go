package fixtures

import (
	"context"
	"testing"
)

type Personality struct {
	ID      *int64 `db:"id"`
	EntryID *int64 `db:"entry_id"`
	Type    string `db:"type"`
}

func NewPersonality(ctx context.Context) *ModelConnector {
	personality := &Personality{
		Type: "introvert",
	}

	return &ModelConnector{
		Model: personality,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Personalities = append(f.Personalities, personality)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				personality.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, personality)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO personality (entry_id, type) VALUES (:entry_id, :type)", personality)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			personality.ID = &id
		},
	}
}
