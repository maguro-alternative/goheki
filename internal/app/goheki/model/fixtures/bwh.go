package fixtures

import (
	"context"
	"testing"
)

type BWH struct {
	EntryID int64  `db:"entry_id"`
	Bust    int64  `db:"bust"`
	Waist   int64  `db:"waist"`
	Hip     int64  `db:"hip"`
	Height  *int64 `db:"height"`
	Weight  *int64 `db:"weight"`
}

func NewBWH(ctx context.Context, setter ...func(b *BWH)) *ModelConnector {
	bwh := &BWH{
		Bust:   1,
		Waist:  1,
		Hip:    1,
		Height: nil,
		Weight: nil,
	}

	//setter(bwh)

	return &ModelConnector{
		Model: bwh,
		setter: func() {
			for _, s := range setter {
				s(bwh)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.BWHs = append(f.BWHs, bwh)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				bwh.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, bwh)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, "INSERT INTO bwh (entry_id, bust, waist, hip, height, weight) VALUES (:entry_id, :bust, :waist, :hip, :height, :weight)", bwh)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
