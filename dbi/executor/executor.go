/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package executor

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Execution is an interface for mocking purposes of sql functions like open(), ping(), exec(), close() etc.
type Execution interface {
	Open(driverName, dataSourceName string) error
	Close() error
	Ping() error
	SwitchToDB(dbName string) error
	Query(name, statement string) (map[string][]interface{}, error)
}

// SQLExecutor keeps handle to sql database and map of prepared queries' statements
type SQLExecutor struct {
	handle *sql.DB
	stmts  map[string]*sql.Stmt
}

// NewExecutor returns a pointer to SQLExecutor with initialized map of stmt
// as an Execution interface defined in that way for mocking purposes
var NewExecutor = func() Execution {
	return &SQLExecutor{stmts: make(map[string]*sql.Stmt)}
}

// Open opens a database specified by its database driver name and a driver-specific
// data source name. To verify that the data source name is valid, call Ping()
// The Open function should be called just once. It is rarely necessary to close a DB.
func (se *SQLExecutor) Open(driverName, dataSourceName string) error {
	var err error
	se.handle, err = sql.Open(driverName, dataSourceName)
	return err
}

// Close closes the database, releasing any open resources. It is rare to Close a DB,
// as the DB handle is meant to be long-lived and shared between many goroutines.
func (se *SQLExecutor) Close() error {
	return se.handle.Close()
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary
func (se *SQLExecutor) Ping() error {
	return se.handle.Ping()
}

// SwitchToDB changes the database context to the specified database
func (se *SQLExecutor) SwitchToDB(dbName string) error {
	_, err := se.handle.Exec("USE " + dbName)
	return err
}

// Query executes a query and returns its output in convenient format (as a map to its values where keys are the names of columns)
func (se *SQLExecutor) Query(name, statement string) (map[string][]interface{}, error) {
	rows, err := execQuery(se, name, statement)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Cannot execute query `%+v`, err=%+v", statement, err)
	}

	// get query output (rows) and parse it to map
	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	if len(cols) == 0 {
		return nil, errors.New("Invalid row does not contain columns")
	}

	table := map[string][]interface{}{}
	vals := make([]interface{}, len(cols))
	valsPtrs := make([]interface{}, len(vals))
	cnt := 0

	for i := range valsPtrs {
		valsPtrs[i] = &vals[i]
	}

	for rows.Next() {
		err = rows.Scan(valsPtrs...)
		if err != nil {
			return nil, err
		}

		for i, val := range vals {
			columnName := strings.ToLower(cols[i])
			table[columnName] = append(table[columnName], val)
		}
		cnt++
	} // end of row.Next()

	return table, nil
}

// execQuery creates a prepared statement and executes a query that returns rows (typically a SELECT statement)
func execQuery(se *SQLExecutor, name, statement string) (*sql.Rows, error) {
	var err error

	// if query statement is not prepared (do not occured in map), prepare it
	if se.stmts[name] == nil {
		// preparing query statement is needed to use the newer protocol for MySQL driver
		// which provides information about type of result's value (can be obtained by using reflection)
		se.stmts[name], err = se.handle.Prepare(statement)
		if err != nil {
			se.stmts[name].Close()
			return nil, err
		}
	}
	// execute query, output data is returned as rows
	return se.stmts[name].Query()
}
