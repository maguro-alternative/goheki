package fixtures

type HairColor struct {
	EntryID *int64 `db:"entry_id"`
	Color   string `db:"color"`
}
