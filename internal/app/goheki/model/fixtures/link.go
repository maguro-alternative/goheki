package fixtures

type Link struct {
	ID       *int64 `db:"id"`
	EntryID  *int64 `db:"entry_id"`
	Type     string `db:"type"`
	URL      string `db:"url"`
	Nsfw     bool   `db:"nsfw"`
	Darkness bool   `db:"darkness"`
}
