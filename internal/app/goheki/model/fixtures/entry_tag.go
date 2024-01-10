package fixtures

type EntryTag struct {
	ID      *int64 `db:"id"`
	EntryID *int64 `db:"entry_id"`
	TagID   *int64 `db:"tag_id"`
}
