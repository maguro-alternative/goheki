package bwh

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type BWH struct {
	EntryID int64  `db:"entry_id" json:"entry_id"`
	Bust    int64  `db:"bust" json:"bust"`
	Waist   int64  `db:"waist" json:"waist"`
	Hip     int64  `db:"hip" json:"hip"`
	Height  *int64 `db:"height" json:"height"`
	Weight  *int64 `db:"weight" json:"weight"`
}

func (b *BWH) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.EntryID, validation.Required),
		validation.Field(&b.Bust, validation.Required),
		validation.Field(&b.Waist, validation.Required),
		validation.Field(&b.Hip, validation.Required),
	)
}

type BWHsJson struct {
	BWHs []BWH `json:"bwhs"`
}

func (b *BWHsJson) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.BWHs, validation.Required),
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
