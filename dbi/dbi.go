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
	"os"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/dtype"
	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/parser"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	// Name of plugin
	Name = "dbi"
	// Version of plugin
	Version = 3
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// DbiPlugin holds information about the configuration database and defined queries
type DbiPlugin struct {
	databases   map[string]*dtype.Database
	queries     map[string]*dtype.Query
	initialized bool
}

// CollectMetrics returns values of desired metrics defined in mts
func (dbiPlg *DbiPlugin) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {

	var err error
	metrics := []plugin.PluginMetricType{}
	data := map[string]interface{}{}

	// initialization - done once
	if dbiPlg.initialized == false {
		// CollectMetrics(mts) is called only when mts has one item at least
		err = dbiPlg.setConfig(mts[0])
		if err != nil {
			// Cannot obtained sql settings
			return nil, err
		}

		err = openDBs(dbiPlg.databases)
		if err != nil {
			return nil, err
		}

		dbiPlg.initialized = true
	} // end of initialization

	// execute dbs queries and get output
	data, err = dbiPlg.executeQueries()
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	for _, m := range mts {
		if value, ok := data[joinNamespace(m.Namespace())]; ok {
			metric := plugin.PluginMetricType{
				Namespace_: m.Namespace(),
				Data_:      value,
				Source_:    hostname,
				Timestamp_: time.Now(),
				Version_:   m.Version(),
			}
			metrics = append(metrics, metric)
		}

	}

	return metrics, nil
}

// GetConfigPolicy returns config policy
func (dbiPlg *DbiPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// GetMetricTypes returns metrics types exposed by snap-plugin-collector-dbi
func (dbiPlg *DbiPlugin) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	metrics := map[string]interface{}{}
	mts := []plugin.PluginMetricType{}

	err := dbiPlg.setConfig(cfg)
	if err != nil {
		// cannot obtained sql settings from Global Config
		return nil, err
	}

	metrics, err = dbiPlg.getMetrics()
	if err != nil {
		return nil, err
	}

	for name := range metrics {
		mts = append(mts, plugin.PluginMetricType{Namespace_: splitNamespace(name)})
	}

	return mts, nil
}

// New returns snap-plugin-collector-dbi instance
func New() *DbiPlugin {
	dbiPlg := &DbiPlugin{databases: map[string]*dtype.Database{}, queries: map[string]*dtype.Query{}, initialized: false}

	return dbiPlg
}

// setConfig extracts config item from Global Config or Metric Config, parses its contents (mainly information
// about databases and queries) and assigned them to appriopriate DBiPlugin fields
func (dbiPlg *DbiPlugin) setConfig(cfg interface{}) error {
	setFile, err := config.GetConfigItem(cfg, "setfile")
	if err != nil {
		// cannot get config item
		return err
	}

	dbiPlg.databases, dbiPlg.queries, err = parser.GetDBItemsFromConfig(setFile.(string))
	if err != nil {
		// cannot parse sql config contents
		return err
	}

	return nil
}

// getMetrics returns map with dbi metrics values, where keys are metrics names
func (dbiPlg *DbiPlugin) getMetrics() (map[string]interface{}, error) {
	metrics := map[string]interface{}{}

	err := openDBs(dbiPlg.databases)
	defer closeDBs(dbiPlg.databases)

	if err != nil {
		return nil, err
	}

	// execute dbs queries and get statement outputs
	metrics, err = dbiPlg.executeQueries()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// executeQueries executes all defined queries of each database and returns results as map to its values,
// where keys are equal to columns' names
func (dbiPlg *DbiPlugin) executeQueries() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	//execute queries for each defined databases
	for dbName, db := range dbiPlg.databases {

		if !db.Active {
			//skip if db is not active (none established connection)
			fmt.Fprintf(os.Stderr, "Cannot execute queries for database %s, is inactive (connection was not established properly)\n", dbName)
			continue
		}

		// retrive name from queries to be executed for this db
		for _, queryName := range db.QrsToExec {
			statement := dbiPlg.queries[queryName].Statement

			out, err := db.Executor.Query(queryName, statement)

			if err != nil {
				// log failing query and take the next one
				fmt.Fprintf(os.Stderr, "Cannot execute query %s for database %s", queryName, dbName)
				continue
			}

			for resName, res := range dbiPlg.queries[queryName].Results {
				instanceOk := false

				// to avoid inconsistency of columns names caused by capital letters (especially for postgresql driver)
				instanceFrom := strings.ToLower(res.InstanceFrom)
				valueFrom := strings.ToLower(res.ValueFrom)

				if !isEmpty(instanceFrom) {
					if len(out[instanceFrom]) == len(out[valueFrom]) {
						instanceOk = true
					}
				}

				for index, value := range out[valueFrom] {
					instance := ""

					if instanceOk {
						instance = fmt.Sprintf("%v", fixDataType(out[instanceFrom][index]))
					}

					key := createNamespace(dbName, resName, res.InstancePrefix, instance)

					if _, exist := data[key]; exist {
						return nil, fmt.Errorf("Namespace `%s` has to be unique, but is not", key)
					}

					data[key] = fixDataType(value)
				}
			}
		} // end of range db_queries_to_execute
	} // end of range databases

	if len(data) == 0 {
		return nil, errors.New("No data obtained from defined queries")
	}

	return data, nil
}

// fixDataType converts `arg` to a string if its type is an array of bytes or time.Time, in other case there is no change
func fixDataType(arg interface{}) interface{} {
	var result interface{}

	switch arg.(type) {
	case []byte:
		result = string(arg.([]byte))

	case time.Time:
		// gob: type time.Time is not registered for gob interface, conversion to string
		result = arg.(time.Time).String()

	default:
		result = arg
	}

	return result
}
