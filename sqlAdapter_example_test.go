package sqladapter_test

import (
	"fmt"

	"os"
	"strconv"

	adapter "github.com/cdleo/go-sql-adapter"
	"github.com/cdleo/go-sql-adapter/engines"
	"github.com/cdleo/go-sql-adapter/loggers"
	"github.com/cdleo/go-sql-adapter/translators"
)

type People struct {
	Id       int    `db:"id"`
	Nombre   string `db:"firstname"`
	Apellido string `db:"lastname"`
}

func Example_sqlAdapter() {

	connector := engines.NewSqlite3Connector(":memory:")
	translator := translators.NewNoopTranslator()
	stdoutLogger, _ := loggers.NewBasicLogger()
	//stdoutLogger.SetLogLevel(logger.LogLevel_Trace.String())
	sqlAdapter := adapter.NewSQLAdapter(connector, translator, stdoutLogger)

	db, err := sqlAdapter.Open()
	if err != nil {
		fmt.Println("Unable to connect to DB")
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		fmt.Printf("Unable to execute statement %v\n", err)
		os.Exit(1)
	}

	stmt, err := db.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = stmt.Exec("Gene", "Kranz")
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	rows, err := db.Query("SELECT id, firstname, lastname FROM people")
	if err != nil {
		fmt.Printf("Unable to query data %v\n", err)
		os.Exit(1)
	}

	var p People
	for rows.Next() {
		_ = rows.Scan(&p.Id, &p.Nombre, &p.Apellido)
		fmt.Println(strconv.Itoa(p.Id) + ": " + p.Nombre + " " + p.Apellido)
	}

	// Output:
	// 1: Gene Kranz
}
