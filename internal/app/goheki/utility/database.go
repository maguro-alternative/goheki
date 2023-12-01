package utility

import (
	"github.com/jmoiron/sqlx"

	"github.com/maguro-alternative/goheki/pkg/db"
)

type botHandlerDB struct {
	DBHandler *db.DBHandler
}

func NewSqlDB(dbSql *sqlx.DB) *botHandlerDB {
	return &botHandlerDB{
		DBHandler: db.NewDBHandler(dbSql),
	}
}