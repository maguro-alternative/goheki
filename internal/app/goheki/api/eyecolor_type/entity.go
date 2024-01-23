package eyecolortype

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type EyeColorType struct {
	ID    int64  `db:"id"`
	Color string `db:"color"`
}

func (e *EyeColorType) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Color, validation.Required),
	)
}

type EyeColorTypesJson struct {
	EyeColorTypes []EyeColorType `json:"eyecolor_types"`
}

func (e *EyeColorTypesJson) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.EyeColorTypes, validation.Required),
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
