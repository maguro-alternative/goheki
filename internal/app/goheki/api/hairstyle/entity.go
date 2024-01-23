package hairstyle

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairStyle struct {
	EntryID int64 `db:"entry_id"`
	StyleID int64 `db:"style_id"`
}

func (h *HairStyle) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.EntryID, validation.Required),
		validation.Field(&h.StyleID, validation.Required),
	)
}

type HairStylesJson struct {
	HairStyles []HairStyle `json:"hair_styles"`
}

func (h *HairStylesJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairStyles, validation.Required),
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
