package eyecolor

type EyeColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

type EyeColorsJson struct {
	EyeColors []EyeColor `json:"eyecolors"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
