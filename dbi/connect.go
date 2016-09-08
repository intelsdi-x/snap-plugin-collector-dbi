// +build linux

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

package dbi

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/dtype"
)

// defaultPort specifies the TCP/IP port on which sql driver is listening
// for connections from client applications
var defaultPort = map[string]string{
	"mysql":    "3306",
	"postgres": "5432",
}

// getDefaultPort returns default port for specific driver
func getDefaultPort(driver string) string {
	return defaultPort[driver]
}

// openDB opens a database and verifies connection by calling ping to it
func openDB(db *dtype.Database) error {
	var dsn string

	// if port is not defined, set defaults
	if isEmpty(db.Port) {
		db.Port = getDefaultPort(db.Driver)
	}

	switch db.Driver {
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			db.Username, db.Password, db.Host, db.Port, db.DBName)

	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			db.Username, db.Password, db.Host, db.Port, db.DBName)

	default:
		return fmt.Errorf("SQL Driver %s is not supported", db.Driver)
	}
	err := db.Executor.Open(db.Driver, dsn)
	if err != nil {
		return err
	}

	// ping db to verify a connection
	if err = db.Executor.Ping(); err != nil {
		return err
	}

	if db.SelectDB != "" {
		// switch the connection when SelectDB is defined in cfg
		err = db.Executor.SwitchToDB(db.SelectDB)
		if err != nil {
			return err
		}
	}

	db.Active = true

	return nil
}

// openDBs opens databases and verifies connections by calling ping to them
func openDBs(dbs map[string]*dtype.Database) error {
	once := false
	for i := range dbs {
		err := openDB(dbs[i])
		if err != nil {
			return err
		}
		once = true
	}

	if !once {
		return errors.New("Cannot open any of defined database")
	}

	return nil
}

// closeDB closes a database
func closeDB(db *dtype.Database) error {
	if db.Active {
		err := db.Executor.Close()
		if err != nil {
			return err
		}
		db.Active = false
	}
	return nil
}

// closeDBs closes databases (exported due to use in main.go)
func closeDBs(dbs map[string]*dtype.Database) []error {
	//errors := []error{}
	var errors []error
	for i := range dbs {
		err := closeDB(dbs[i])
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
