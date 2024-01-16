package fixtures

import (
	"context"
	"testing"
)

type HairColorType struct {
	ID    *int64 `db:"id"`
	Color string `db:"color"`
}

func NewHairColorType(ctx context.Context, setter ...func(h *HairColorType)) *ModelConnector {
	hairColorType := &HairColorType{
		Color: "black",
	}

	//setter(hairColor)

	return &ModelConnector{
		Model: hairColorType,
		setter: func() {
			for _, s := range setter {
				s(hairColorType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairColorTypes = append(f.HairColorTypes, hairColorType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *HairColor:
				hairColor := connectingModel.(*HairColor)
				hairColor.ColorID = *hairColorType.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hairColorType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			r := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO haircolor_type (
					color
				) VALUES (
					$1
				) RETURNING id`,
				hairColorType.Color,
			).Scan(&hairColorType.ID)
			if r != nil {
				t.Fatalf("insert error: %v", r)
			}
		},
	}
}
