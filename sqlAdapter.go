package sqladapter

import (
	"database/sql"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-sql-adapter/loggers"
)

type sqlAdapter struct {
	connector SQLEngineConnector
}

func NewSQLAdapter(connector SQLEngineConnector, translator SQLSyntaxTranslator, logger logger.Logger) *sqlAdapter {

	if logger == nil {
		logger = loggers.NewNullLogger()
	}

	Register(connector, translator, logger)

	return &sqlAdapter{
		connector,
	}
}

func (s *sqlAdapter) Open() (*sql.DB, error) {

	if db, err := s.connector.Open(); err != nil {
		return nil, s.connector.ErrorHandler(err)
	} else {
		return db, nil
	}
}
