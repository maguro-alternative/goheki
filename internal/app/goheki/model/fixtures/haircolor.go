package fixtures

import (
	"context"
	"testing"
)

type HairColor struct {
	EntryID int64 `db:"entry_id"`
	Color   string `db:"color"`
}

func NewHairColor(ctx context.Context) *ModelConnector {
	hairColor := &HairColor{
		Color: "black",
	}

	return &ModelConnector{
		Model: hairColor,
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairColors = append(f.HairColors, hairColor)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				hairColor.EntryID = *entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hairColor)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.dbv1.NamedExecContext(ctx, "INSERT INTO hair_color (entry_id, color) VALUES (:entry_id, :color)", hairColor)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}