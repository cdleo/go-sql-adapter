package sqladapter

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/cdleo/go-commons/logger"
	proxy "github.com/shogo82148/go-sql-proxy"
)

func Register(connector SQLEngineConnector, translator SQLSyntaxTranslator, logger logger.Logger) {

	if isAlreadyRegistered(connector.DriverName()) {
		return
	}

	sql.Register(connector.DriverName(), proxy.NewProxyContext(connector.Driver(), &proxy.HooksContext{
		Open: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			logger.Qry("Open")
			return nil
		},
		Begin: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			logger.Qry("Begin")
			return nil
		},
		Commit: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			logger.Qry("Commit")
			return nil
		},
		Rollback: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			logger.Qry("Rollback")
			return nil
		},
		PrePrepare: func(_ context.Context, stmt *proxy.Stmt) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return nil, nil
		},
		PreQuery: func(_ context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		Query: func(_ context.Context, _ interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows) error {
			logger.Qryf("Query: %s; args = %v", prettyQuery(stmt.QueryString), args)
			return nil
		},
		PostQuery: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows, err error) error {
			if err != nil {
				return connector.ErrorHandler(err)
			}
			logger.Tracef("Query: %s; args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
			return nil
		},
		PreExec: func(_ context.Context, stmt *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		Exec: func(_ context.Context, _ interface{}, stmt *proxy.Stmt, args []driver.NamedValue, result driver.Result) error {
			logger.Qryf("Exec: %s; args = %v", prettyQuery(stmt.QueryString), args)
			return nil
		},
		PostExec: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, err error) error {
			if err != nil {
				return connector.ErrorHandler(err)
			}
			logger.Tracef("Exec: %s; args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
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
