package fixtures

import (
	"context"
	"testing"
)

type HairColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

func NewHairColor(ctx context.Context, setter ...func(h *HairColor)) *ModelConnector {
	hairColor := &HairColor{}

	//setter(hairColor)

	return &ModelConnector{
		Model: hairColor,
		setter: func() {
			for _, s := range setter {
				s(hairColor)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairColors = append(f.HairColors, hairColor)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				hairColor.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hairColor)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO haircolor (entry_id, color_id) VALUES (:entry_id, :color_id)", hairColor)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
