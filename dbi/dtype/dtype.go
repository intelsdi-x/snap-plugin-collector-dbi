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

package dtype

import (
	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/executor"
)

// Database holds connection information (driver, host, username etc.),
// names of queries to perform and instance of executor which stores handle to db
type Database struct {
	Driver    string
	Host      string
	Port      string
	Username  string
	Password  string
	DBName    string
	SelectDB  string
	Executor  executor.Execution
	Active    bool
	QrsToExec []string // names of queries to be executed for the database
}

// Query holds statement of the query and its results (there is one or more) which
// structure defines how the returned data should be interpreted
type Query struct {
	Statement string
	Results   map[string]Result
}

// Result holds information specified the columns whose values will be used to
// distinguish results defined by `InstanceFrom` (additionally prefix can be added)
// or whose content will be used as the actual data dfined by `ValueFrom.
type Result struct {
	InstanceFrom   string
	InstancePrefix string
	ValueFrom      string
}
