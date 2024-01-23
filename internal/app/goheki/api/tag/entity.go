package tag

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func (t *Tag) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Name, validation.Required),
	)
}

type TagsJson struct {
	Tags []Tag `json:"tags"`
}

func (t *TagsJson) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Tags, validation.Required),
	)
}

type IDs struct {
	IDs []int64 `json:"ids"`
}

func (i *IDs) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.IDs, validation.Required),
	)
}
