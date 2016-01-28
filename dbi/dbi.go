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
	"regexp"
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
	Version = 1
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

		// prepare namespace to regular expression
		name := joinNamespace(m.Namespace())
		name = strings.Replace(name, "/", "\\/", -1)
		name = strings.Replace(name, "*", ".*", -1)
		regex := regexp.MustCompile("^" + name + "$")

		for key := range data {
			match := regex.FindStringSubmatch(key)

			if match == nil {
				continue
			}

			if value, ok := data[key]; ok {
				metric := plugin.PluginMetricType{
					Namespace_: splitNamespace(key),
					Data_:      value,
					Source_:    hostname,
					Timestamp_: time.Now(),
					Version_:   m.Version(),
				}
				metrics = append(metrics, metric)
			}
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

	// add supporting of whitecards
	mts = append(mts, plugin.PluginMetricType{Namespace_: append(nsPrefix, "*")})

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
				if !isEmpty(res.InstanceFrom) {
					if len(out[res.InstanceFrom]) == len(out[res.ValueFrom]) {
						instanceOk = true
					}
				}

				for index, value := range out[res.ValueFrom] {
					instance := ""
					if instanceOk {
						instance = fmt.Sprintf("%s", out[res.InstanceFrom][index])
					}

					key := createNamespace(dbName, queryName, resName, res.InstancePrefix, instance)
					data[key] = value
				}
			}
		} // end of range db_queries_to_execute
	} // end of range databases

	if len(data) == 0 {
		return nil, errors.New("No data obtained from defined queries")
	}

	return data, nil
}
