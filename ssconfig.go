// MIT License
//
// Copyright (c) 2021 Josh Simonot
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Super Simple Config
package ssconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// FieldError identifies a field which failed to load
// and its associated error
type FieldError struct {
	Field string
	Error error
}

// ConfigError lists all fields that failed to load
type ConfigError struct {
	Fields []FieldError
}

// Error lists all config errors
func (err ConfigError) Error() string {
	return fmt.Sprintf("ConfigErrors: %+v", err.Fields)
}

// Set config options
type Set struct {
	EnvPrefix string
	FilePath  string
}

// Load with default (nil) options
func Load(conf interface{}) *ConfigError {
	return Set{}.Load(conf)
}

// Load with Set config options
func (ssc Set) Load(conf interface{}) *ConfigError {
	var confError ConfigError

	// load values from file if it exists
	if confFile, err := ioutil.ReadFile(ssc.FilePath); err == nil {

		err = json.Unmarshal(confFile, &conf)
		if err != nil {
			confError.Fields = append(confError.Fields, FieldError{ssc.FilePath, err})
			log.Printf("ssconfig: failed to parse file: %s\n", ssc.FilePath)
		} else {
			log.Printf("ssconfig: %s ✓", ssc.FilePath)
		}

	} else if ssc.FilePath != "" {
		confError.Fields = append(confError.Fields, FieldError{ssc.FilePath, err})
		log.Printf("ssconfig: %s (file not found)\n", ssc.FilePath)
	}

	// load values from ENV if they exist, overwritting
	// the values from file. Only works for struct fields
	confValue := reflect.ValueOf(conf).Elem()
	confType := confValue.Type()

	if confType.Kind() == reflect.Struct {
		for i := 0; i < confType.NumField(); i++ {

			field := confValue.Field(i)
			confName := confType.Field(i).Name

			if field.CanSet() {
				env := os.Getenv(ssc.EnvPrefix + confName)
				if env != "" {

					switch field.Kind() {

					case reflect.String:
						// string is handled separately since
						// json.Unmarshal expects quotes around
						// string values
						field.SetString(env)

					default:
						err := json.Unmarshal([]byte(env), field.Addr().Interface())
						if err != nil {
							log.Printf("ssconfig: %s (type %s): %s\n", confName, field.Type().String(), err)
							confError.Fields = append(confError.Fields, FieldError{confName, err})
							continue
						}
					}
					log.Printf("ssconfig: %s ✓", ssc.EnvPrefix+confName)
				}
			}
		}
	}

	if len(confError.Fields) > 0 {
		return &confError
	}
	return nil
}
