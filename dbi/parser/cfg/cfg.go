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

package cfg

// To unmarshal JSON into a struct, structs have to contain exported fields

type SQLConfig struct {
	Queries   []QueryType     `json:"queries"`
	Databases []DatabasesType `json:"databases"`
}

type QueryType struct {
	Name      string            `json:"name"`
	Statement string            `json:"statement"`
	Results   []QueryResultType `json:"results"`
}

type QueryResultType struct {
	ResultName     string `json:"name"`
	InstanceFrom   string `json:"instance_from"`
	InstancePrefix string `json:"instance_prefix"`
	ValueFrom      string `json:"value_from"`
}

type DatabasesType struct {
	Name           string           `json:"name"`
	Driver         string           `json:"driver"`
	DriverOption   DriverOptionType `json:"driver_option"`
	SelectDb       string           `json:"selectdb"`
	QueryToExecute []DBQueryType    `json:"dbqueries"`
}

type DBQueryType struct {
	QueryName string `json:"query"`
}

type DriverOptionType struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	Port     string `json:"port"`
}
