package haircolor

type HairColor struct {
	EntryID int64 `db:"entry_id"`
	ColorID int64 `db:"color_id"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
