package link

type Link struct {
	ID       *int64 `db:"id"`
	EntryID  *int64 `db:"entry_id"`
	Type     string `db:"type"`
	URL      string `db:"url"`
	Nsfw     bool   `db:"nsfw"`
	Darkness bool   `db:"darkness"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
