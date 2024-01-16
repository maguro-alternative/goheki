package hairstyle

type HairStyle struct {
	EntryID int64 `db:"entry_id"`
	StyleID int64 `db:"style_id"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
