package service

import (
	"github.com/jmoiron/sqlx"

	"github.com/maguro-alternative/goheki/pkg/db"
)

type databaseHandler struct {
	DBHandler *db.DBHandler
}

func NewSqlDB(dbSql *sqlx.DB) *databaseHandler {
	return &databaseHandler{
		DBHandler: db.NewDBHandler(dbSql),
	}
}