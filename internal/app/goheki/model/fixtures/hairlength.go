package fixtures

type hairLength struct {
	EntryID *int64 `db:"entry_id"`
	Length  int64  `db:"length"`
}
