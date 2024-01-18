package personalitytype

type PersonalityType struct {
	ID   *int64 `db:"id"`
	Type string `db:"type"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
