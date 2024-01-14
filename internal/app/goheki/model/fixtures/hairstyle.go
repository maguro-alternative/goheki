package fixtures

import (
	"context"
	"testing"
)

type HairStyle struct {
	EntryID *int64 `db:"entry_id"`
	Style   string `db:"style"`
}

func NewHairStyle(ctx context.Context, setter func(h *HairStyle)) *ModelConnector {
	hairStyle := &HairStyle{
		Style: "long",
	}

	setter(hairStyle)

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
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hairstyle (entry_id, style) VALUES (:entry_id, :style)", hairStyle)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
