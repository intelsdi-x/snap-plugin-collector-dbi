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
	"github.com/intelsdi-x/snap/core"
)

var (

	// Mts is a mocked metrics
	Mts = []plugin.MetricType{

		// 1) For db = dbName1 one query=q1 is executed which has no additional info (like name or prefix)
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName1", "categoryA")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName1", "categoryB")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName1", "categoryC")},

		// 2) For db = dbName2 two queries are executed:
		// a) query=q1 - has no additional info (like name or prefix)
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "categoryA")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "categoryB")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "categoryC")},

		// b) query=q2 - has defined additional info (like resultName or instancePrefix) which are appended to a namespace
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryA")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryB")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName1", "category_prefix", "categoryC")},

		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName2", "categoryA")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName2", "categoryB")},
		plugin.MetricType{Namespace_: core.NewNamespace("intel", "dbi", "dbName2", "rName2", "categoryC")},
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
	FileName = "temp_setfile.json"

	SetfileCorr   = "mock/corrMockSetfile.json"
	SetfileIncorr = "mock/incorrMockSetfile.json"
)
