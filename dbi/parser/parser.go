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

package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/dtype"
	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/executor"
	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/parser/cfg"
)

// Parser holds maps to queries and databases
type Parser struct {
	qrs map[string]*dtype.Query
	dbs map[string]*dtype.Database
}

// GetDBItemsFromConfig parses the contents of the file `fName` and returns maps to
// databases and queries instances which structurs are pre-defined in package dtype
func GetDBItemsFromConfig(fName string) (map[string]*dtype.Database, map[string]*dtype.Query, error) {

	var sqlCnf cfg.SQLConfig

	if strings.ContainsAny(fName, "$") {
		// filename contains environment variable, expand it
		fName = expandFileName(fName)
	}

	data, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 0 {
		return nil, nil, fmt.Errorf("SQL settings file `%v` is empty", fName)
	}

	err = json.Unmarshal(data, &sqlCnf)

	if err != nil {
		return nil, nil, fmt.Errorf("Invalid structure of file `%v` to be unmarshalled", fName)
	}

	p := &Parser{
		qrs: map[string]*dtype.Query{},
		dbs: map[string]*dtype.Database{},
	}

	for _, query := range sqlCnf.Queries {
		err := p.addQuery(query)
		if err != nil {
			return nil, nil, err
		}

	}

	for _, db := range sqlCnf.Databases {
		err := p.addDatabase(db)
		if err != nil {
			return nil, nil, err
		}
	}

	return p.dbs, p.qrs, nil
}

// addDatabase adds database instance to databases
func (p *Parser) addDatabase(dt cfg.DatabasesType) error {

	if len(strings.TrimSpace(dt.Name)) == 0 {
		return fmt.Errorf("Data name is empty")
	}

	if _, exist := p.dbs[dt.Name]; exist {
		return fmt.Errorf("Data name `%+s` is not unique", dt.Name)
	}

	//getting info about which queries are to be executed
	execQrs := []string{}
	for _, q := range dt.QueryToExecute {
		execQrs = append(execQrs, q.QueryName)
	}

	// adding database to databases map
	p.dbs[dt.Name] = &dtype.Database{
		Driver:    dt.Driver,
		Host:      dt.DriverOption.Host,
		Port:      dt.DriverOption.Port,
		Username:  dt.DriverOption.Username,
		Password:  dt.DriverOption.Password,
		DBName:    dt.DriverOption.DbName,
		SelectDB:  dt.SelectDb,
		Active:    false,
		QrsToExec: execQrs,
		Executor:  executor.NewExecutor(),
	}

	return nil
}

// addQuery adds query instance to queries
func (p *Parser) addQuery(qt cfg.QueryType) error {

	if len(strings.TrimSpace(qt.Name)) == 0 {
		return fmt.Errorf("Query name is empty")
	}

	if _, exist := p.qrs[qt.Name]; exist {
		return fmt.Errorf("Query name `%+s` is not unique", qt.Name)
	}

	results := map[string]dtype.Result{}

	for _, r := range qt.Results {

		if _, exist := results[r.ResultName]; exist {
			return fmt.Errorf("Query `%+s` has result `%+s` which name is not unique", qt.Name, r.ResultName)
		}

		// add result to the map `results`
		results[r.ResultName] = dtype.Result{
			InstanceFrom:   r.InstanceFrom,
			InstancePrefix: r.InstancePrefix,
			ValueFrom:      r.ValueFrom,
		}

	} // end of range q.Results

	// adding query to queries map
	p.qrs[qt.Name] = &dtype.Query{
		Statement: qt.Statement,
		Results:   results,
	}
	return nil
}

// expandFileName replaces name of environment variable with its value and returns expanded filename
func expandFileName(fName string) string {

	// split namespace to get its components
	fNameCmps := strings.Split(fName, "/")

	for i, fNameCmp := range fNameCmps {
		if strings.Contains(fNameCmp, "$") {
			envName := strings.TrimPrefix(fNameCmp, "$")
			if envValue := os.Getenv(envName); envValue != "" {
				// replace name of environment variable with its value
				fNameCmps[i] = envValue
			}
		}
	}
	return strings.Join(fNameCmps, "/")
}
