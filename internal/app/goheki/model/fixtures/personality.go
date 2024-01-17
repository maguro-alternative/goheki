package fixtures

import (
	"context"
	"testing"
)

type Personality struct {
	EntryID int64 `db:"entry_id"`
	TypeID  int64 `db:"type_id"`
}

func NewPersonality(ctx context.Context, setter ...func(p *Personality)) *ModelConnector {
	personality := &Personality{}

	//setter(personality)

	return &ModelConnector{
		Model: personality,
		setter: func() {
			for _, s := range setter {
				s(personality)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Personalities = append(f.Personalities, personality)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			case *Entry:
				entry := connectingModel.(*Entry)
				personality.EntryID = *entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, personality)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(
				ctx,
				`INSERT INTO personality (
					entry_id,
					type_id
				) VALUES (
					:entry_id,
					:type_id
				)`,
				personality,
			)
			if err != nil {
				t.Fatalf("insert error: %v %v", err, personality)
			}
		},
	}
}
