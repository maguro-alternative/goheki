package eyecolor

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type EyeColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

func (e *EyeColor) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.EntryID, validation.Required),
		validation.Field(&e.ColorID, validation.Required),
	)
}

type EyeColorsJson struct {
	EyeColors []EyeColor `json:"eyecolors"`
}

func (e *EyeColorsJson) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.EyeColors, validation.Required),
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
