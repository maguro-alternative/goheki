package hairlengthtype

type HairLengthType struct {
	ID     *int64 `db:"id"`
	Length string `db:"length"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
