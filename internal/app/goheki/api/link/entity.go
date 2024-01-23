package link

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Link struct {
	ID       int64 `db:"id"`
	EntryID  int64 `db:"entry_id"`
	Type     string `db:"type"`
	URL      string `db:"url"`
	Nsfw     bool   `db:"nsfw"`
	Darkness bool   `db:"darkness"`
}

func (l *Link) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.EntryID, validation.Required),
		validation.Field(&l.Type, validation.Required),
		validation.Field(&l.URL, validation.Required),
		validation.Field(&l.Nsfw, validation.Required),
		validation.Field(&l.Darkness, validation.Required),
	)
}

type LinksJson struct {
	Links []Link `json:"links"`
}

func (l *LinksJson) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Links, validation.Required),
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
