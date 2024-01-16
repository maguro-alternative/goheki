package fixtures

import (
	"context"
	"testing"
)

type Personality struct {
	EntryID *int64 `db:"entry_id"`
	Type    string `db:"type"`
}

func NewPersonality(ctx context.Context, setter ...func(p *Personality)) *ModelConnector {
	personality := &Personality{
		Type: "introvert",
	}

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
				personality.EntryID = entry.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, personality)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO personality (
					entry_id,
					type
				) VALUES (
					$1,
					$2
				)`,
				personality.EntryID,
				personality.Type,
			)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
