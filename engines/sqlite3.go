package engines

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	adapter "github.com/cdleo/go-sql-adapter"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/mattn/go-sqlite3"
)

type sqlite3Conn struct {
	url string
}

const sqlite3_DriverName = "sqlite3-adapter"

func NewSqlite3Connector(url string) adapter.SQLEngineConnector {
	return &sqlite3Conn{
		url,
	}
}

func (s *sqlite3Conn) Open() (*sql.DB, error) {
	return sql.Open(sqlite3_DriverName, s.url)
}

func (s *sqlite3Conn) DriverName() string {
	return sqlite3_DriverName
}

func (s *sqlite3Conn) Driver() driver.Driver {
	return &sqlite3.SQLiteDriver{}
}

func (s *sqlite3Conn) ErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	if sqliteError, ok := err.(sqlite3.Error); ok {

		if sqliteError.Code == 18 { //SQLITE_TOOBIG
			return sqlcommons.ValueTooLargeForColumn

		} else if sqliteError.Code == 19 { //SQLITE_CONSTRAINT
			if sqliteError.ExtendedCode == 787 || /*SQLITE_CONSTRAINT_FOREIGNKEY*/
				sqliteError.ExtendedCode == 1555 || /*SQLITE_CONSTRAINT_PRIMARYKEY*/
				sqliteError.ExtendedCode == 1811 /*SQLITE_CONSTRAINT_TRIGGER*/ {
				return sqlcommons.IntegrityConstraintViolation

			} else if sqliteError.ExtendedCode == 1299 { //SQLITE_CONSTRAINT_NOTNULL
				return sqlcommons.CannotSetNullColumn

			} else if sqliteError.ExtendedCode == 2067 { //SQLITE_CONSTRAINT_UNIQUE
				return sqlcommons.UniqueConstraintViolation

			}
		} else if sqliteError.Code == 25 { //SQLITE_RANGE
			return sqlcommons.InvalidNumericValue
		}

		return fmt.Errorf("Unhandled SQLite3 error. Code:[%s] Extended:[%s] Desc:[%s]", sqliteError.Code, sqliteError.ExtendedCode, sqliteError.Error())

	} else {
		return err
	}
}
