package hairlength

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairLength struct {
	EntryID          int64 `db:"entry_id"`
	HairLengthTypeID int64 `db:"hairlength_type_id"`
}

func (h *HairLength) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.EntryID, validation.Required),
		validation.Field(&h.HairLengthTypeID, validation.Required),
	)
}

type HairLengthsJson struct {
	HairLengths []HairLength `json:"hairlengths"`
}

func (h *HairLengthsJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairLengths, validation.Required),
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
