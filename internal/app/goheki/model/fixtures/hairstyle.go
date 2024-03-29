package fixtures

import (
	"context"
	"testing"
)

type HairStyle struct {
	EntryID int64 `db:"entry_id"`
	StyleID int64 `db:"style_id"`
}

func NewHairStyle(ctx context.Context, setter ...func(h *HairStyle)) *ModelConnector {
	hairStyle := &HairStyle{}

	//setter(hairStyle)

	return &ModelConnector{
		Model: hairStyle,
		setter: func() {
			for _, s := range setter {
				s(hairStyle)
			}
		},
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
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hairstyle (entry_id, style_id) VALUES (:entry_id, :style_id)", hairStyle)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
