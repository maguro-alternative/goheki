package personalitytype

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type PersonalityType struct {
	ID   int64  `db:"id"`
	Type string `db:"type"`
}

func (p *PersonalityType) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, validation.Required),
	)
}

type PersonalityTypesJson struct {
	PersonalityTypes []PersonalityType `json:"personality_types"`
}

func (p *PersonalityTypesJson) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.PersonalityTypes, validation.Required),
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
