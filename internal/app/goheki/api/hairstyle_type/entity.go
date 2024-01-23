package hairstyletype

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairStyleType struct {
	ID    int64  `db:"id"`
	Style string `db:"style"`
}

func (h *HairStyleType) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Style, validation.Required),
	)
}

type HairStyleTypesJson struct {
	HairStyleTypes []HairStyleType `json:"hairstyle_types"`
}

func (h *HairStyleTypesJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairStyleTypes, validation.Required),
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
