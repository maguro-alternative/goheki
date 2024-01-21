package eyecolortype

type EyeColorType struct {
	ID    *int64 `db:"id"`
	Color string `db:"color"`
}

type EyeColorTypesJson struct {
	EyeColorTypes []EyeColorType `json:"eyecolor_types"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
