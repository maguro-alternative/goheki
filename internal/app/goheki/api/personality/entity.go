package personality

type Personality struct {
	EntryID *int64 `db:"entry_id"`
	TypeID  *int64 `db:"type_id"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
