package fixtures

type BWH struct {
	EntryID *int64 `db:"entry_id"`
	Bust	int64 `db:"bust"`
	Waist	int64 `db:"waist"`
	Hip		int64 `db:"hip"`
	Height	*int64 `db:"height"`
	Weight	*int64 `db:"weight"`
}