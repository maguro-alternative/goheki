package fixtures

import (
	"context"
	"testing"
)

type HairLengthType struct {
	ID     *int64 `db:"id"`
	Length string `db:"length"`
}

func NewHairLengthType(ctx context.Context, setter ...func(h *HairLengthType)) *ModelConnector {
	heirLengthType := &HairLengthType{}

	//setter(heirLength)

	return &ModelConnector{
		Model: heirLengthType,
		setter: func() {
			for _, s := range setter {
				s(heirLengthType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.HairLengthTypes = append(f.HairLengthTypes, heirLengthType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, heirLengthType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			r, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO hairlength_type (length) VALUES (:length) RETURNING id", heirLengthType)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
			heirLengthType.ID = &id
		},
	}
}
