package fixtures

import (
	"context"
	"testing"
)

type HairStyle struct {
	EntryID *int64 `db:"entry_id"`
	Style   string `db:"style"`
}

func NewHairStyle(ctx context.Context) *ModelConnector {
	hairStyle := &HairStyle{
		Style: "long",
	}

	return &ModelConnector{
		Model: hairStyle,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairStyles = append(f.HairStyles, hairStyle)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				hairStyle.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hairStyle)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hair_style (entry_id, style) VALUES (:entry_id, :style)", hairStyle)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			// 連番されるIDを取得する
			id, err := result.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			// 連番されるIDをセットする
			hairStyle.EntryID = &id
		},
	}
}
