package bwh

type BWH struct {
	EntryID *int64 `db:"entry_id" json:"entry_id"`
	Bust    int64  `db:"bust" json:"bust"`
	Waist   int64  `db:"waist" json:"waist"`
	Hip     int64  `db:"hip" json:"hip"`
	Height  *int64 `db:"height" json:"height"`
	Weight  *int64 `db:"weight" json:"weight"`
}

type BWHsJson struct {
	BWHs []BWH `json:"bwhs"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
