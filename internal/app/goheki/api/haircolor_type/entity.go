package haircolortype

type HairColorType struct {
	ID    *int64 `db:"id"`
	Color string `db:"color"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
