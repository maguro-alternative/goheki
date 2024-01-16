package fixtures

import (
	"context"
	"testing"
)

type HairStyleType struct {
	ID    *int64 `db:"id"`
	Style string `db:"style"`
}

func NewHairStyleType(ctx context.Context, setter ...func(h *HairStyleType)) *ModelConnector {
	hairStyleType := &HairStyleType{
		Style: "long",
	}

	//setter(hairStyle)

	return &ModelConnector{
		Model: hairStyleType,
		setter: func() {
			for _, s := range setter {
				s(hairStyleType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairStyleTypes = append(f.HairStyleTypes, hairStyleType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, hairStyleType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			r := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO hairstyle_type (style) VALUES ($1) RETURNING id`,
				hairStyleType.Style,
			).Scan(&hairStyleType.ID)
			if r != nil {
				t.Fatalf("insert error: %v", r)
			}
		},
	}
}
