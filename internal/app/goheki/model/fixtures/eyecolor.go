package fixtures

import (
	"context"
	"testing"
)

type EyeColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

func NewEyeColor(ctx context.Context, setter ...func(e *EyeColor)) *ModelConnector {
	eyeColor := &EyeColor{}

	//setter(eyeColor)

	return &ModelConnector{
		Model: eyeColor,
		setter: func() {
			for _, s := range setter {
				s(eyeColor)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.EyeColors = append(f.EyeColors, eyeColor)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				eyeColor.EntryID = entry.ID
			case *EyeColorType:
				eyeColorType := connectingModel.(*EyeColorType)
				eyeColor.ColorID = eyeColorType.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, eyeColor)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(
				ctx,
				`INSERT INTO eyecolor (
					entry_id,
					color_id
				) VALUES (
					:entry_id,
					:color_id
				)`,
				eyeColor,
			)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
