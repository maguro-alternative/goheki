package fixtures

import (
	"context"
	"testing"
)

type HairLength struct {
	EntryID *int64 `db:"entry_id"`
	Length  int64  `db:"length"`
}

func NewHeirLength(ctx context.Context) *ModelConnector {
	heirLength := &HairLength{
		Length: 1,
	}

	return &ModelConnector{
		Model: heirLength,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairLengths = append(f.HairLengths, heirLength)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				heirLength.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, heirLength)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hair_length (entry_id, length) VALUES (:entry_id, :length)", heirLength)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			heirLength.EntryID = &id
		},
	}
}
