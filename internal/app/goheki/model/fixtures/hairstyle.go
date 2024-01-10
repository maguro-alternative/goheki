package fixtures

type HairStyle struct {
	EntryID *int64 `db:"entry_id"`
	Style   string `db:"style"`
}
