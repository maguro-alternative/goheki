package fixtures

import (
	"context"
	"testing"
)

type EyeColorType struct {
	ID    *int64 `db:"id"`
	Color string `db:"color"`
}

func NewEyeColorType(ctx context.Context, setter ...func(e *EyeColorType)) *ModelConnector {
	eyeColorType := &EyeColorType{
		Color: "black",
	}

	//setter(eyeColor)

	return &ModelConnector{
		Model: eyeColorType,
		setter: func() {
			for _, s := range setter {
				s(eyeColorType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.EyeColorTypes = append(f.EyeColorTypes, eyeColorType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *EyeColor:
				eyeColor := connectingModel.(*EyeColor)
				eyeColor.ColorID = *eyeColorType.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, eyeColorType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			r := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO eyecolor_type (
					color
				) VALUES (
					$1
				) RETURNING id`,
				eyeColorType.Color,
			).Scan(&eyeColorType.ID)
			if r != nil {
				t.Fatalf("insert error: %v", r)
			}
		},
	}
}
