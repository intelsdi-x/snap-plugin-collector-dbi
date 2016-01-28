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

import "strings"

// nsPrefix is prefix of metrics namespace
var nsPrefix = []string{"intel", "dbi"}

// notAllowedChars contains all not allowed chars in namespace
var notAllowedChars = []string{" ", "-", "(", ")", "[", "]", "{", "}", ",", ";"}

// joinNamespace concatenates the elements of namespace to create a single string separated by slash
func joinNamespace(ns []string) (name string) {
	return "/" + strings.Join(ns, "/")
}

// splitNamespace splits name and returns a slice of the substrings between slashes
func splitNamespace(name string) (ns []string) {
	return strings.Split(strings.TrimPrefix(name, "/"), "/")
}

// createNamespace returns metric namespace
func createNamespace(dbName, queryName, resultName, instancePrefix, instanceValue string) string {

	ns := append(nsPrefix, dbName, queryName)

	//append resultName (omit if empty or equal to queryName)
	if isNotEmpty(resultName) && (queryName != resultName) {
		ns = append(ns, resultName)
	}

	// append instancePrefix (omit if empty)
	if isNotEmpty(instancePrefix) {
		ns = append(ns, instancePrefix)
	}

	// append instanceValue (omit if empty)
	if isNotEmpty(instanceValue) {
		ns = append(ns, instanceValue)
	}

	return validateNamespace(joinNamespace(ns))
}

// validateNamespace removes not allowed chars from namespace
func validateNamespace(str string) string {

	// replace notAllowedChars to underscore
	for _, c := range notAllowedChars {
		str = strings.Replace(str, c, "_", -1)

		// to avoid double undescores
		str = strings.Replace(str, "__", "_", -1)
	}

	// trimming white space
	str = strings.TrimSpace(str)

	return str
}

// isEmpty returns true when string `str` is empty
func isEmpty(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}

// isEmpty returns true when string `str` is not empty
func isNotEmpty(str string) bool {
	return !isEmpty(str)
}
