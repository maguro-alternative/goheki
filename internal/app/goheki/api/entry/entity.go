package entry

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Entry struct {
	ID        int64     `db:"id" json:"id"`
	SourceID  int64     `db:"source_id" json:"source_id"`
	Name      string    `db:"name" json:"name"`
	Image     string    `db:"image" json:"image"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (e *Entry) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.SourceID, validation.Required),
		validation.Field(&e.Name, validation.Required),
		validation.Field(&e.Image, validation.Required),
		validation.Field(&e.Content, validation.Required),
		validation.Field(&e.CreatedAt, validation.Required),
	)
}

type EntriesJson struct {
	Entries []Entry `json:"entries"`
}

func (e *EntriesJson) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Entries, validation.Required),
	)
}

type Source struct {
	ID   *int64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Url  string `db:"url" json:"url"`
	Type string `db:"type" json:"type"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}

func (i *IDs) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.IDs, validation.Required),
	)
}
