package hairstyletype

type HairStyleType struct {
	ID    *int64 `db:"id"`
	Style string `db:"style"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
