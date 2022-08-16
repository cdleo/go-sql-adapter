package sqladapter_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/cdleo/go-commons/sqlcommons"
	adapter "github.com/cdleo/go-sql-adapter"
	"github.com/cdleo/go-sql-adapter/engines"
	enginesMocks "github.com/cdleo/go-sql-adapter/engines/mocks"
	"github.com/cdleo/go-sql-adapter/translators"
	"github.com/stretchr/testify/require"
)

type Customers struct {
	Id         int           `db:"id"`
	Name       string        `db:"name"`
	Updatetime time.Time     `db:"updatetime"`
	Age        sql.NullInt64 `db:"age"`
	Group      int           `db:"cust_group"`
	Dummy      string        `db:"not_existing_field"`
}

func Test_sqlConn_InitErr(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(false)
	translator := translators.NewNoopTranslator()

	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)
	_, err := sqlConn.Open()
	require.Error(t, err)
}

func Test_sqlConn_InitOK(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(true)
	translator := translators.NewNoopTranslator()

	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)
	_, err := sqlConn.Open()
	require.NoError(t, err)
}

func Test_sqlConn_CreateTables(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	// Exec
	require.NoError(t, createTablesHelper(db))

	db.Close()
}

func Test_sqlConn_DropTables(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)
	require.NoError(t, createTablesHelper(db))

	// Exec
	require.NoError(t, dropTablesHelper(db))

	db.Close()
}

func Test_sqlConn_StoreData(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))

	// Exec
	require.NoError(t, insertDataHelper(db))

	db.Close()
}

func Test_sqlConn_ReturnData(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	rows, err2 := db.Query("SELECT name FROM customers")
	defer rows.Close()

	require.NoError(t, err2)
	require.True(t, rows.Next())

	db.Close()
}

func Test_sqlConn_CanThrowInvalidTableError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	_, err2 := db.Query("SELECT name FROM customerxs")

	require.Error(t, err2)

	db.Close()
}

func Test_sqlConn_CanThrowCannotInsertNullError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))

	// Exec
	_, err2 := db.Exec("INSERT INTO customers (name, updatetime) VALUES (:1,:2)", nil, time.Now())

	require.ErrorIs(t, err2, sqlcommons.CannotSetNullColumn)

	db.Close()
}

func Test_sqlConn_CanThrowCannotUpdateNullError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	_, err2 := db.Exec("UPDATE customers c SET name = :1 WHERE c.name = :2", nil, "Pablo")

	require.Error(t, err2, sqlcommons.CannotSetNullColumn)

	db.Close()
}

func Test_sqlConn_CanThrowUniqueConstraintViolationError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:?_foreign_keys=on")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	_, err2 := db.Exec("INSERT INTO customers (name, updatetime, age, cust_group)VALUES(:1, :2, :3, :4)", "Juan", time.Now(), nil, 1)

	require.ErrorIs(t, err2, sqlcommons.UniqueConstraintViolation)

	db.Close()
}

func Test_sqlConn_CanThrowForeignKeyConstraintViolationError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:?_foreign_keys=on")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	_, err2 := db.Exec("UPDATE customers SET cust_group = :1 WHERE name = :2", 2, "Pablo")

	require.ErrorIs(t, err2, sqlcommons.IntegrityConstraintViolation)

	db.Close()
}

func Test_sqlConn_CanThrowIntegrityConstraintViolationError(t *testing.T) {
	// Setup
	connector := engines.NewSqlite3Connector(":memory:?_foreign_keys=on")
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(db))
	require.NoError(t, insertDataHelper(db))

	// Exec
	_, err2 := db.Exec("DELETE from customers_groups WHERE id = :1", 1)

	require.ErrorIs(t, err2, sqlcommons.IntegrityConstraintViolation)

	db.Close()
}

func Test_sqlConn_CanThrowValueTooLargeError(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(true)
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	// Exec
	query := fmt.Sprintf("INSERT INTO customers (name, cust_group) VALUES (%s,%d)", "'verylongname'", 1)
	connector.PatchExec(query, sqlcommons.ValueTooLargeForColumn)

	_, err2 := db.Exec(query)

	require.ErrorIs(t, err2, sqlcommons.ValueTooLargeForColumn)

	db.Close()
}

func Test_sqlConn_CanThrowSubqueryReturnsMoreThanOneRowError(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(true)
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	query := "SELECT name FROM customers WHERE id = (SELECT id FROM customers)"
	connector.PatchQuery(query, nil, nil, sqlcommons.SubqueryReturnsMoreThanOneRow)

	// Exec
	_, err2 := db.Query(query)

	require.ErrorIs(t, err2, sqlcommons.SubqueryReturnsMoreThanOneRow)

	db.Close()
}

func Test_sqlConn_CanThrowInvalidNumericValueError(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(true)
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	query := "UPDATE customers SET age = :1 WHERE name = :2"
	connector.PatchExec(query, sqlcommons.InvalidNumericValue, "twelve", "Pablo")

	// Exec
	_, err2 := db.Exec(query, "twelve", "Pablo")

	require.ErrorIs(t, err2, sqlcommons.InvalidNumericValue)

	db.Close()
}

func Test_sqlConn_CanThrowValueLargerThanPrecisionError(t *testing.T) {
	// Setup
	connector := enginesMocks.NewMockSQLConnector(true)
	translator := translators.NewNoopTranslator()
	sqlConn := adapter.NewSQLAdapter(connector, translator, nil)

	db, err := sqlConn.Open()
	require.NoError(t, err)

	query := "UPDATE customers SET age = :1 WHERE name = :2"
	connector.PatchExec(query, sqlcommons.ValueLargerThanPrecision, 949.0044, "Pablo")

	// Exec
	_, err2 := db.Exec(query, 949.0044, "Pablo")

	require.ErrorIs(t, err2, sqlcommons.ValueLargerThanPrecision)

	db.Close()
}

func dropTablesHelper(db *sql.DB) error {

	if _, err := db.Exec(`DROP TABLE IF EXISTS customers_groups`); err != nil {
		return err
	}
	if _, err := db.Exec(`DROP TABLE IF EXISTS customers`); err != nil {
		return err
	}
	return nil
}

func createTablesHelper(db *sql.DB) error {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS customers_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		groupname TEXT NOT NULL)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name CHAR(10) NOT NULL,
		updatetime TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
		age INT NULL,
		cust_group INT NOT NULL,
		FOREIGN KEY (cust_group) REFERENCES customers_groups (id) ON DELETE RESTRICT
		CONSTRAINT customers_un UNIQUE (name))`); err != nil {
		return err
	}
	return nil
}

func insertDataHelper(db *sql.DB) error {

	if _, err := db.Exec(`INSERT INTO customers_groups (groupname) VALUES('General');`); err != nil {
		return err
	}

	if statement, err := db.Prepare("INSERT INTO customers (name, updatetime, age, cust_group)VALUES(:1, :2, :3, :4)"); err != nil {
		return err
	} else {
		if _, err := statement.Exec("Juan", time.Now(), nil, 1); err != nil {
			return err
		}
		if _, err := statement.Exec("Pedro", time.Now(), nil, 1); err != nil {
			return err
		}
		if _, err := statement.Exec("Pablo", time.Now(), 99, 1); err != nil {
			return err
		}
	}

	return nil
}
