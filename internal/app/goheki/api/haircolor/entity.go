package haircolor

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type HairColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

func (h *HairColor) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.EntryID, validation.Required),
		validation.Field(&h.ColorID, validation.Required),
	)
}

type HairColorsJson struct {
	HairColors []HairColor `json:"haircolors"`
}

func (h *HairColorsJson) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.HairColors, validation.Required),
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
