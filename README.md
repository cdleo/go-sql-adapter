# go-sql-adapter

[![Go Reference](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/cdleo/go-sql-adapter) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/cdleo/go-sql-adapter/master/LICENSE) [![Build Status](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/badges/build.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/build-status/main) [![Code Coverage](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/badges/coverage.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/?branch=main) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sql-adapter/?branch=main)


## General

**go-sql-adapter** it's a muti DB Engine adapter over the GO **database/sql** package. It provides a set of standard error codes, providing abstraction from the implemented DB engine, allowing change the DB just modifying the configuration and without need adjust or modify the source code.
Besides that, provides a very limited cross-engine sql translator.

Another feature included in the package it's the automatic logging of the executed sentences and it's elapsed time, trougth the implementation of the Logger interface: [github.com/cdleo/go-commons/logger/logger.go](https://github.com/cdleo/go-commons/logger/logger.go)

This package provides a nolog and a basic implementation of that interface, but if you need a more customized and powerful implementation, you can use the following one: [github.com/cdleo/go-zla](https://github.com/cdleo/go-zla)


**Supported Engines**
Currently, the next set of engines are supported:
- **Oracle**: Using the godror driver [github.com/godror/godror](https://github.com/godror/godror)
- **Postgres**: Using the pq driver [github.com/lib/pq](https://github.com/lib/pq)
- **SQLite3**: Using the go-sqlite3 driver [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)


**Usage**
This example program shows the initialization and the use at basic level:
```go
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
	stdoutLogger := loggers.NewBasicLogger()
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
```

## Sample

You can find a sample of the use of go-sql-adapter project [HERE](https://github.com/cdleo/go-sql-adapter/blob/master/sqlAdapter_example_test.go)

## Contributing

Comments, suggestions and/or recommendations are always welcomed. Please check the [Contributing Guide](CONTRIBUTING.md) to learn how to get started contributing.
