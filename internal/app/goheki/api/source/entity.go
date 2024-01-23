package source

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Source struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Url  string `db:"url" json:"url"`
	Type string `db:"type" json:"type"`
}

func (s *Source) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Url, validation.Required),
		validation.Field(&s.Type, validation.Required),
	)
}

type SourcesJson struct {
	Sources []Source `json:"sources"`
}

func (s *SourcesJson) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Sources, validation.Required),
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
