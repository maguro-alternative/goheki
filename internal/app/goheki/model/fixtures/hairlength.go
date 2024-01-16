package fixtures

import (
	"context"
	"testing"
)

type HairLength struct {
	EntryID          int64 `db:"entry_id"`
	HairLengthTypeID int64 `db:"hairlength_type_id"`
}

func NewHairLength(ctx context.Context, setter ...func(h *HairLength)) *ModelConnector {
	heirLength := &HairLength{}

	//setter(heirLength)

	return &ModelConnector{
		Model: heirLength,
		setter: func() {
			for _, s := range setter {
				s(heirLength)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairLengths = append(f.HairLengths, heirLength)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				heirLength.EntryID = *entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, heirLength)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hairlength (entry_id, hairlength_type_id) VALUES (:entry_id, :hairlength_type_id)", heirLength)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
