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

package mockdata

import (
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
)

var (

	// Mts is a mocked metrics
	Mts = []plugin.PluginMetricType{

		// 1) For db = dbName1 one query=q1 is executed which has no additional info (like name or prefix)
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName1", "categoryA"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName1", "categoryB"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName1", "categoryC"}},

		// 2) For db = dbName2 two queries are executed:
		// a) query=q1 - has no additional info (like name or prefix)
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "categoryA"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "categoryB"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "categoryC"}},

		// b) query=q2 - has defined additional info (like resultName or instancePrefix) which are appended to a namespace
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryA"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryB"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryC"}},

		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName2", "categoryA"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName2", "categoryB"}},
		plugin.PluginMetricType{Namespace_: []string{"intel", "dbi", "dbName2", "rName2", "categoryC"}},
	}

	// QueryOutput is a mocked query output
	QueryOutput = map[string][]interface{}{
		"category": []interface{}{[]byte(`categoryA`), []byte(`categoryB`), "categoryC"},
		"value":    []interface{}{-10.5, 0.0, 10.5},
	}

	// QueryOutputTimestamp is a mocked query output, where values are timestamps
	QueryOutputTimestamp = map[string][]interface{}{
		"category": []interface{}{[]byte(`categoryA`), []byte(`categoryB`), "categoryC"},
		"value":    []interface{}{time.Now(), time.Now().Add(1 * time.Hour), time.Now().Add(2 * time.Hour)},
	}

	// FileName is a path of mock setfile
	FileName = "./temp_setfile.json"

	// FileCont is a mocked content of setfile
	FileCont = []byte(`{
		    "queries": [
		        {
		            "name": "q1",
		            "statement": "statementA",
		            "results": [
		                {   "name": "",                
		                    "instance_from": "category",
		                    "value_from": "value"
		                }
		            ]
		        },
		        {
		            "name": "q2",
		            "statement": "statementB",
		            "results": [
		                {
		                    "name": "rName1",
		                    "instance_from": "category",
		                    "instance_prefix": "category prefix",
		                    "value_from": "value"
		                },
		                {
		                    "name": "rName2",
		                    "instance_from": "category",
		                    "value_from": "value"
		                }
		            ]
		        }
		    ],
		    "databases": [
		        {
		            "name": "dbName1",
		            "driver": "mysql",
		            "driver_option": {
		                "host": "localhost",
		                "port": "3306",
		                "username": "tester",
		                "password": "passwd",
		                "dbname": "mydb"
		            },
		            "dbqueries": [
		                {
		                    "query": "q1"
		                }
		            ]
		        },
		        {
		            "name": "dbName2",
		            "driver": "postgres",
		            "driver_option": {
		                "host": "localhost",
		                "username": "tester",
		                "password": "passwd",
		                "dbname": "mydb"
		            },
		            "selectdb": "slctdb",
		            "dbqueries": [
						{
		                    "query": "q1"                
						},
		                {
		                    "query": "q2"
		                }
		            ]
		        },
				
				{
		            "name": "db3",
		            "driver": "unknown",
		            "driver_option": {
		                "host": "localhost",
		                "username": "tester",
		                "password": "passwd",
		                "dbname": "mydb"
		            },
		
		            "dbqueries": [
		                {
		                    "query": "q1"
		                }
		            ]
		        }
		    ]
		}`)
)
