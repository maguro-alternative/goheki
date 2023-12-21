package service

import (
	"github.com/maguro-alternative/goheki/configs/envconfig"

	"github.com/jmoiron/sqlx"
	"github.com/gorilla/sessions"
)

type IndexService struct {
	DB             *sqlx.DB
	CookieStore    *sessions.CookieStore
	Env            *envconfig.Env
}

// NewTODOService returns new TODOService.
func NewIndexService(
	db *sqlx.DB,
	cookieStore *sessions.CookieStore,
	env *envconfig.Env,
) *IndexService {
	return &IndexService{
		DB:             db,
		CookieStore:    cookieStore,
		Env:            env,
	}
}