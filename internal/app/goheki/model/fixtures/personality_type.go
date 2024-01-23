package fixtures

import (
	"context"
	"testing"
)

type PersonalityType struct {
	ID      int64 `db:"id"`
	Type    string `db:"type"`
}

func NewPersonalityType(ctx context.Context, setter ...func(p *PersonalityType)) *ModelConnector {
	personalityType := &PersonalityType{
		Type: "introvert",
	}

	//setter(personality)

	return &ModelConnector{
		Model: personalityType,
		setter: func() {
			for _, s := range setter {
				s(personalityType)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.PersonalityTypes = append(f.PersonalityTypes, personalityType)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, personalityType)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			// 連番されるIDをセットする
			result := f.DBv1.QueryRowxContext(
				ctx,
				`INSERT INTO personality_type (
					type
				) VALUES (
					$1
				) RETURNING id`,
				personalityType.Type,
			).Scan(&personalityType.ID)
			if result != nil {
				t.Fatalf("insert error: %v", result)
			}
		},
	}
}
