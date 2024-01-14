package hairlength

type HairLength struct {
	EntryID *int64 `db:"entry_id"`
	Length  string  `db:"length"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
