package source

type Source struct {
	ID      *int64 `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Url     string `db:"url" json:"url"`
	Type    string `db:"type" json:"type"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
