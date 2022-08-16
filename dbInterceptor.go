package sqladapter

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/cdleo/go-commons/logger"
	proxy "github.com/cdleo/go-sql-proxy"
)

func Register(connector SQLEngineConnector, translator SQLSyntaxTranslator, logger logger.Logger) {

	if isAlreadyRegistered(connector.DriverName()) {
		return
	}

	sql.Register(connector.DriverName(), proxy.NewProxyContext(connector.Driver(), &proxy.HooksContext{
		PostOpen: func(_ context.Context, _ interface{}, conn *proxy.Conn, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Open (source: %s)", err.Error())
				return stdErr
			}
			logger.Qry("Open")
			return nil
		},
		PostBegin: func(_ context.Context, _ interface{}, conn *proxy.Conn, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Begin (source: %s)", err.Error())
				return stdErr
			}
			logger.Qry("Begin")
			return nil
		},
		PostCommit: func(_ context.Context, _ interface{}, tx *proxy.Tx, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Commit (source: %s)", err.Error())
				return stdErr
			}
			logger.Qry("Commit")
			return nil
		},
		PostRollback: func(_ context.Context, _ interface{}, tx *proxy.Tx, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Rollback (source: %s)", err.Error())
				return stdErr
			}
			logger.Qry("Rollback")
			return nil
		},
		PrePrepare: func(_ context.Context, stmt *proxy.Stmt) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return nil, nil
		},
		PostPrepare: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Prepare (source: %s)", err.Error())
				return stdErr
			}
			logger.Tracef("Prepare: %s", prettyQuery(stmt.QueryString))
			return nil
		},
		PreQuery: func(_ context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		PostQuery: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Query: %s; Args = %v (source: %s)", prettyQuery(stmt.QueryString), args, err.Error())
				return stdErr
			}
			logger.Tracef("Query: %s; Args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
			return nil
		},
		PreExec: func(_ context.Context, stmt *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		PostExec: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, err error) error {
			if err != nil {
				stdErr := connector.ErrorHandler(err)
				logger.Errorf(stdErr, "Exec: %s; Args = %v (source: %s)", prettyQuery(stmt.QueryString), args, err.Error())
				return stdErr
			}
			logger.Tracef("Exec: %s; Args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
			return nil
		},
	}))
}

func prettyQuery(query string) string {
	return strings.ReplaceAll(strings.ReplaceAll(query, "\t", ""), "\n", "")
}

func isAlreadyRegistered(proxyName string) bool {

	drivers := sql.Drivers()
	for _, item := range drivers {
		if item == proxyName {
			return true
		}
	}

	return false
}
