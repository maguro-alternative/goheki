package entry

import (
	"time"
)

type Entry struct {
	ID        *int64    `db:"id" json:"id"`
	SourceID  int64     `db:"source_id" json:"source_id"`
	Name      string    `db:"name" json:"name"`
	Image     string    `db:"image" json:"image"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Source struct {
	ID   *int64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Url  string `db:"url" json:"url"`
	Type string `db:"type" json:"type"`
}

type ID struct {
	ID int64 `json:"id"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}
