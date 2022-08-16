package sqladapter

import (
	"database/sql"
	"database/sql/driver"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-sql-adapter/loggers"
)

type MockSQLEngineConnector interface {
	SQLEngineConnector

	PatchBegin(err error)
	PatchCommit(err error)
	PatchRollback(err error)

	PatchExec(query string, err error, args ...driver.Value)
	PatchQuery(query string, columns []string, values []driver.Value, err error, args ...driver.Value)
	PatchQueryRow(query string, result map[string]string, err error)
}

type SQLEngineConnector interface {
	Open() (*sql.DB, error)
	DriverName() string
	Driver() driver.Driver
	ErrorHandler(err error) error
}

type SQLSyntax int

const (
	SQLSyntax_Oracle SQLSyntax = iota
	SQLSyntax_PostgreSQL
	SQLSyntax_SQLite3
)

//go:generate mockgen -package translatorsMocks -destination translators/mocks/sqlSyntaxTranslator.go . SQLSyntaxTranslator
type SQLSyntaxTranslator interface {
	Translate(query string) string
}

type sqlDB struct {
	connector  SQLEngineConnector
	translator SQLSyntaxTranslator
	logger     logger.Logger
}

func NewSQLAdapter(connector SQLEngineConnector, translator SQLSyntaxTranslator, logger logger.Logger) *sqlDB {

	if logger == nil {
		logger, _ = loggers.NewNullLogger()
	}

	return &sqlDB{
		connector,
		translator,
		logger,
	}
}

func (s *sqlDB) Open() (*sql.DB, error) {

	Register(s.connector, s.translator, s.logger)

	if db, err := s.connector.Open(); err != nil {
		return nil, s.connector.ErrorHandler(err)
	} else {
		return db, nil
	}
}
