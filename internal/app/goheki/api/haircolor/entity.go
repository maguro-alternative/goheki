package haircolor

type HairColor struct {
	EntryID int64  `db:"entry_id"`
	Color   string `db:"color"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
