package personality

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Personality struct {
	EntryID int64 `db:"entry_id"`
	TypeID  int64 `db:"type_id"`
}

func (p *Personality) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.EntryID, validation.Required),
		validation.Field(&p.TypeID, validation.Required),
	)
}

type PersonalitiesJson struct {
	Personalities []Personality `json:"personalities"`
}

func (p *PersonalitiesJson) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Personalities, validation.Required),
	)
}

type IDs struct {
	IDs []int64 `json:"ids"`
}

func (i *IDs) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.IDs, validation.Required),
	)
}
