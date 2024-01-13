package personality

type Personality struct {
	ID      *int64 `db:"id"`
	EntryID *int64 `db:"entry_id"`
	Type    string `db:"type"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
