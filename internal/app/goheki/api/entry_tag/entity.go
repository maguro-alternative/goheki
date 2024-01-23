package entry_tag

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type EntryTag struct {
	ID      int64 `db:"id" json:"id"`
	EntryID int64 `db:"entry_id" json:"entry_id"`
	TagID   int64 `db:"tag_id" json:"tag_id"`
}

func (e *EntryTag) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.EntryID, validation.Required),
		validation.Field(&e.TagID, validation.Required),
	)
}

type EntryTagsJson struct {
	EntryTags []EntryTag `json:"entry_tags"`
}

func (e *EntryTagsJson) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.EntryTags, validation.Required),
	)
}

type Source struct {
	ID      int64 `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Url     string `db:"url" json:"url"`
	Type    string `db:"type" json:"type"`
}

type Entry struct {
	ID        int64    `db:"id" json:"id"`
	SourceID  int64     `db:"source_id" json:"source_id"`
	Name      string    `db:"name" json:"name"`
	Image     string    `db:"image" json:"image"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Tag struct {
	ID        int64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
}

type IDs struct {
	IDs []int64 `json:"ids"`
}

func (i *IDs) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.IDs, validation.Required),
	)
}
