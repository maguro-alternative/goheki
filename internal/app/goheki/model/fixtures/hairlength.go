package fixtures

import (
	"context"
	"testing"
)

type HairLength struct {
	EntryID *int64 `db:"entry_id"`
	Length  int64  `db:"length"`
}

func NewHeirLength(ctx context.Context, setter func(h *HairLength)) *ModelConnector {
	heirLength := &HairLength{
		Length: 1,
	}

	setter(heirLength)

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
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hairlength (entry_id, length) VALUES (:entry_id, :length)", heirLength)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
