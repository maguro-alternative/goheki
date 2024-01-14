package hekiradarchart

type HekiRadarChart struct {
	EntryID *int64 `db:"entry_id"`
	AI      int64  `db:"ai"`
	NU      int64  `db:"nu"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
