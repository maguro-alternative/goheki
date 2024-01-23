package hairlengthtype

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairLengthType struct {
	ID     int64 `db:"id"`
	Length string `db:"length"`
}

func (h *HairLengthType) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Length, validation.Required),
	)
}

type HairLengthTypesJson struct {
	HairLengthTypes []HairLengthType `json:"hairlength_types"`
}

func (h *HairLengthTypesJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairLengthTypes, validation.Required),
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
