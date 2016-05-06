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
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/executor"
	"github.com/intelsdi-x/snap-plugin-collector-dbi/dbi/mock"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type mcMock struct {
	mock.Mock

	handle *sql.DB
	stmts  map[string]*sql.Stmt
}

func (mc *mcMock) Open(driverName, dataSourceName string) error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *mcMock) Close() error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *mcMock) Ping() error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *mcMock) SwitchToDB(dbName string) error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *mcMock) Query(name, statement string) (map[string][]interface{}, error) {
	args := mc.Called()
	return args.Get(0).(map[string][]interface{}), args.Error(1)
}

// mockExecution mocks outputs of Execution SQL methods like Open(), Ping(), Close(), Query() etc.
func (mc *mcMock) mockExecution(errOpen, errClose, errPing, errSwitchToDB, errQuery error, outQuery map[string][]interface{}) {
	mc.On("Open").Return(errOpen)
	mc.On("Close").Return(errClose)
	mc.On("Ping").Return(errPing)
	mc.On("SwitchToDB").Return(errSwitchToDB)
	mc.On("Query").Return(outQuery, errQuery)

	// mock NewExecutor() from `executor` package
	executor.NewExecutor = func() executor.Execution {
		return mc
	}
}

func TestGetConfigPolicy(t *testing.T) {
	dbiPlugin := New()

	Convey("getting config policy", t, func() {
		So(func() { dbiPlugin.GetConfigPolicy() }, ShouldNotPanic)
		_, err := dbiPlugin.GetConfigPolicy()
		So(err, ShouldBeNil)
	})
}

func TestGetMetricTypes(t *testing.T) {

	Convey("getting exposed metric types", t, func() {

		Convey("when no configuration item available", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeNil)
		})

		Convey("when path to setfile is incorrect", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			// mockdata.FileName has not existed yet
			deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeNil)
		})

		Convey("when setfile is empty", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			os.Create(mockdata.FileName)
			defer os.Remove(mockdata.FileName)
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeNil)
		})

		Convey("when cannot open db", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				errors.New("x"), // errOpen
				nil,             // errClose
				nil,             // errPing
				nil,             // errSwitchToDB
				nil,             // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when cannot open, neither close db", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				errors.New("x"), // errOpen
				errors.New("x"), // errClose
				nil,             // errPing
				nil,             // errSwitchToDB
				nil,             // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when cannot ping the open db", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				nil,             // errOpen
				nil,             // errClose
				errors.New("x"), // errPing
				nil,             // errSwitchToDB
				nil,             // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when cannot switch to selected db", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				nil,             // errOpen
				nil,             // errClose
				nil,             // errPing
				errors.New("x"), // errSwitchToDB
				nil,             // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when execution of query returns error", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				nil,                        // errOpen
				nil,                        // errClose
				nil,                        // errPing
				nil,                        // errSwitchToDB
				errors.New("x"),            // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when query returns empty output", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})

			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				nil, // errOpen
				nil, // errClose
				nil, // errPing
				nil, // errSwitchToDB
				nil, // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)

		})

		Convey("successfully obtain metrics name", func() {
			cfg := plugin.NewPluginConfigType()
			dbiPlugin := New()
			createMockFile()
			defer deleteMockFile()
			cfg.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})

			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			// prepare mock output data of query execution
			mockQueryOut := mockdata.QueryOutput

			//mockExecution outputs
			mc.mockExecution(nil, nil, nil, nil, nil, mockQueryOut)

			So(func() { dbiPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := dbiPlugin.GetMetricTypes(cfg)
			So(err, ShouldBeNil)
			So(results, ShouldNotBeEmpty)
		})

	})
}

func TestCollectMetrics(t *testing.T) {

	Convey("when no configuration settings available", t, func() {
		dbiPlugin := New()
		mts := mockdata.Mts
		config := cdata.NewNode()

		for i, _ := range mts {
			mts[i].Config_ = config
		}

		So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
		results, err := dbiPlugin.CollectMetrics(mts)
		So(err, ShouldNotBeNil)
		So(results, ShouldBeEmpty)
	})

	Convey("when configuration is settings are invalid", t, func() {
		mts := mockdata.Mts

		// set metrics config
		config := cdata.NewNode()
		config.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
		for i, _ := range mts {
			mts[i].Config_ = config
		}

		Convey("incorrect path to setfile", func() {
			// mockdata.FileName has not existed (remove it just in case)
			deleteMockFile()
			dbiPlugin := New()
			So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
			results, err := dbiPlugin.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("setfile is empty", func() {
			dbiPlugin := New()
			os.Create(mockdata.FileName)
			defer os.Remove(mockdata.FileName)
			So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
			results, err := dbiPlugin.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

	})

	Convey("SQL methods returns error", t, func() {

		mts := mockdata.Mts
		config := cdata.NewNode()
		config.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
		for i, _ := range mts {
			mts[i].Config_ = config
		}

		createMockFile()
		defer deleteMockFile()

		Convey("when cannot connect to databases", func() {
			dbiPlugin := New()
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				errors.New("x"), // errOpen
				nil,             // errClose
				nil,             // errPing
				nil,             // errSwitchToDB
				nil,             // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
			results, err := dbiPlugin.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when execution of query returns error", func() {
			dbiPlugin := New()
			mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

			//mockExecution outputs
			mc.mockExecution(
				nil,                        // errOpen
				nil,                        // errClose
				nil,                        // errPing
				nil,                        // errSwitchToDB
				errors.New("x"),            // errQuery
				map[string][]interface{}{}, // outQuery
			)

			So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
			results, err := dbiPlugin.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

	})

	Convey("when metric's value is a timestamp - special usecase", t, func() {
		// to cover func fixDataType for time.Time case
		createMockFile()
		defer deleteMockFile()

		dbiPlugin := New()
		mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

		//mockExecution outputs
		mc.mockExecution(
			nil, // errOpen
			nil, // errClose
			nil, // errPing
			nil, // errSwitchToDB
			nil, // errQuery
			mockdata.QueryOutputTimestamp, // outQuery
		)

		mts := mockdata.Mts
		config := cdata.NewNode()
		config.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
		for i, _ := range mts {
			mts[i].Config_ = config
		}

		So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
		results, err := dbiPlugin.CollectMetrics(mts)
		So(err, ShouldBeNil)
		So(len(results), ShouldEqual, len(mts))
	})

	Convey("collect metrics successfully", t, func() {
		createMockFile()
		defer deleteMockFile()

		dbiPlugin := New()
		mc := &mcMock{stmts: make(map[string]*sql.Stmt)}

		//mockExecution outputs
		mc.mockExecution(
			nil,                  // errOpen
			nil,                  // errClose
			nil,                  // errPing
			nil,                  // errSwitchToDB
			nil,                  // errQuery
			mockdata.QueryOutput, // outQuery
		)

		mts := mockdata.Mts
		config := cdata.NewNode()
		config.AddItem("setfile", ctypes.ConfigValueStr{Value: mockdata.FileName})
		mts[0].Config_ = config

		So(func() { dbiPlugin.CollectMetrics(mts) }, ShouldNotPanic)
		results, err := dbiPlugin.CollectMetrics(mts)
		So(err, ShouldBeNil)
		So(len(results), ShouldEqual, len(mts))

	})

}

func createMockFile() {
	deleteMockFile()

	f, _ := os.Create(mockdata.FileName)
	f.Write(mockdata.FileCont)
}

func deleteMockFile() {
	os.Remove(mockdata.FileName)
}
