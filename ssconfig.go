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
	"strconv"
)

// ConfigError lists all fields that failed to load
type ConfigError struct {
	Fields []string
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
func Load(conf interface{}) error {
	return Set{}.Load(conf)
}

// Load with Set config options
func (ssc Set) Load(conf interface{}) error {
	var confError ConfigError

	// load values from file if it exists
	if confFile, err := ioutil.ReadFile(ssc.FilePath); err == nil {
		log.Printf("ssconfig: %s\n", ssc.FilePath)
		err = json.Unmarshal(confFile, &conf)

		if err != nil {
			confError.Fields = append(confError.Fields, fmt.Sprintf("failed to parse file: %s", err.Error()))
			log.Printf("ssconfig: failed to parse file: %s\n", ssc.FilePath)
		}

	} else if ssc.FilePath != "" {
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

					log.Printf("ssconfig: %s\n", ssc.EnvPrefix+confName)
					switch field.Kind() {

					case reflect.String:
						field.SetString(env)

					case reflect.Int8:
						fallthrough

					case reflect.Int16:
						fallthrough

					case reflect.Int32:
						fallthrough

					case reflect.Int64:
						fallthrough

					case reflect.Int:
						if ival, err := strconv.ParseInt(env, 10, 64); err == nil && !field.OverflowInt(ival) {
							field.SetInt(ival)
						} else {
							log.Printf("ssconfig: failed to parse env %s as int64\n", confName)
							confError.Fields = append(confError.Fields, confName)
						}

					case reflect.Float32:
						fallthrough

					case reflect.Float64:
						if fval, err := strconv.ParseFloat(env, 64); err == nil && !field.OverflowFloat(fval) {
							field.SetFloat(fval)
						} else {
							log.Printf("ssconfig: failed to parse env %s as float64\n", confName)
							confError.Fields = append(confError.Fields, confName)
						}

					default:
						log.Printf("ssconfig: %s type not supported by env: %+v\n", confName, field.Kind())
						confError.Fields = append(confError.Fields, fmt.Sprintf("%s: type not supported", confName))
					}
				}
			}
		}
	}

	if len(confError.Fields) > 0 {
		return confError
	}
	return nil
}
