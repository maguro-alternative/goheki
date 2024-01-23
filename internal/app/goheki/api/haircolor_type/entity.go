package haircolortype

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairColorType struct {
	ID    int64 `db:"id"`
	Color string `db:"color"`
}

func (h *HairColorType) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Color, validation.Required),
	)
}

type HairColorTypesJson struct {
	HairColorTypes []HairColorType `json:"haircolor_types"`
}

func (h *HairColorTypesJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairColorTypes, validation.Required),
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
