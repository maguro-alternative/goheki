package hairstyle

type HairStyle struct {
	EntryID *int64 `db:"entry_id"`
	Style   string `db:"style"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
