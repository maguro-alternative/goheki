package hairlength

type HairLength struct {
	EntryID          int64 `db:"entry_id"`
	HairLengthTypeID int64 `db:"hairlength_type_id"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
