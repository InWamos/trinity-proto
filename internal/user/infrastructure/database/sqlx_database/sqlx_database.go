package sqlxdatabase

import "github.com/jmoiron/sqlx"

type SQLXDatabase struct {
	engine *sqlx.E
}
