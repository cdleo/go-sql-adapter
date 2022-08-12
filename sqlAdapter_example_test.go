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

	connector := engines.NewSqlite3Adapter(":memory:")
	translator := translators.NewNoopTranslator()
	logger, _ := loggers.NewBasicLogger()
	sqlAdapter := adapter.NewSQLAdapter(connector, translator, logger)

	db, err := sqlAdapter.Open()
	if err != nil {
		fmt.Println("Unable to connect to DB")
		os.Exit(1)
	}
	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	statement, err = db.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = db.Exec("Gene", "Kranz")
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
