package entrytag

type EntryTag struct {
	ID      *int64 `db:"id" json:"id"`
	EntryID *int64 `db:"entry_id" json:"entry_id"`
	TagID   *int64 `db:"tag_id" json:"tag_id"`
}
