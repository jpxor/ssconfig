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

// Super Simple Config Test
package ssconfig

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestSuccessPaths(t *testing.T) {
	os.Setenv("EnvVarInt", "42")
	os.Setenv("EnvVarInt8", "42")
	os.Setenv("EnvVarInt16", "42")
	os.Setenv("EnvVarInt32", "42")
	os.Setenv("EnvVarInt64", "42")
	os.Setenv("EnvVarString", "42")
	os.Setenv("EnvVarFloat32", "4.2")
	os.Setenv("EnvVarFloat64", "4.2")

	os.Setenv("OverwriteInt", "42")
	os.Setenv("OverwriteInt8", "42")
	os.Setenv("OverwriteInt16", "42")
	os.Setenv("OverwriteInt32", "42")
	os.Setenv("OverwriteInt64", "42")
	os.Setenv("OverwriteString", "42")
	os.Setenv("OverwriteFloat32", "4.2")
	os.Setenv("OverwriteFloat64", "4.2")

	type testestconfig struct {
		// File only
		FileInt    int
		FileString string
		FileList   []float32
		FileMap    map[string]int

		// Env var only
		EnvVarInt     int
		EnvVarInt8    int8
		EnvVarInt16   int16
		EnvVarInt32   int32
		EnvVarInt64   int64
		EnvVarString  string
		EnvVarFloat32 float32
		EnvVarFloat64 float64

		// Env var overwrite file
		OverwriteInt     int
		OverwriteInt8    int8
		OverwriteInt16   int16
		OverwriteInt32   int32
		OverwriteInt64   int64
		OverwriteString  string
		OverwriteFloat32 float32
		OverwriteFloat64 float64
	}

	var testconf testestconfig
	err := Set{FilePath: "./test.json"}.Load(&testconf)

	if err != nil {
		t.Errorf("got %+v, expected no errors", err.Error())
	}

	answerconf := testestconfig{
		FileInt:    42,
		FileString: "42",
		FileList:   []float32{4.2, 4.2, 4.2},
		FileMap:    map[string]int{"answer": 42},

		EnvVarInt:     42,
		EnvVarInt8:    42,
		EnvVarInt16:   42,
		EnvVarInt32:   42,
		EnvVarInt64:   42,
		EnvVarString:  "42",
		EnvVarFloat32: 4.2,
		EnvVarFloat64: 4.2,

		OverwriteInt:     42,
		OverwriteInt8:    42,
		OverwriteInt16:   42,
		OverwriteInt32:   42,
		OverwriteInt64:   42,
		OverwriteString:  "42",
		OverwriteFloat32: 4.2,
		OverwriteFloat64: 4.2,
	}

	if !reflect.DeepEqual(testconf, answerconf) {
		testValue := reflect.ValueOf(testconf)
		answerValue := reflect.ValueOf(answerconf)
		confType := testValue.Type()

		for i := 0; i < confType.NumField(); i++ {
			t.Run(confType.Field(i).Name, func(t *testing.T) {

				if !reflect.DeepEqual(testValue.Field(i).Interface(), answerValue.Field(i).Interface()) {
					t.Errorf("got %+v, want %+v", testValue.Field(i), answerValue.Field(i))
				}
			})
		}
	}
}

func TestDefaultLoad(t *testing.T) {
	os.Setenv("FromEnv", "42")
	type testconf struct {
		FromEnv string
	}
	var tconf testconf
	Load(&tconf)

	if tconf.FromEnv != "42" {
		t.Errorf("got %+v, want %+v", tconf.FromEnv, "42")
	}
}

func TestErrorFileNotFound(t *testing.T) {
	os.Setenv("FromEnv", "42")
	type testconf struct {
		FromEnv string
	}
	var tconf testconf
	err := Set{FilePath: "no/file.conf"}.Load(&tconf)

	// check for expected error
	if err == nil {
		t.Error("got no error, expected file not found")
		return
	}
	if len(err.Fields) != 1 {
		t.Errorf("got %d errors, expected 1 error", len(err.Fields))
		return
	}
	var expectedErr *os.PathError
	testErr := err.Fields[0].Error

	if !errors.As(testErr, &expectedErr) {
		t.Errorf("got %s, expected os.PathError", testErr.Error())
		return
	}

	// ensure load from env vars still succeeds
	if tconf.FromEnv != "42" {
		t.Errorf("got %+v, want %+v", tconf.FromEnv, "42")
	}
}

func TestErrorFileSyntax(t *testing.T) {
	os.Setenv("FromEnv", "42")

	type testconf struct {
		FromEnv string
	}
	var tconf testconf
	err := Set{FilePath: "test.yml"}.Load(&tconf)

	// check for expected error
	if err == nil {
		t.Error("got no error, expected syntax error")
		return
	}
	if len(err.Fields) != 1 {
		t.Errorf("got %d errors, expected 1 error", len(err.Fields))
		return
	}
	var expectedErr *json.SyntaxError
	testErr := err.Fields[0].Error

	if !errors.As(testErr, &expectedErr) {
		t.Errorf("got %s, expected json.SyntaxError", testErr.Error())
		return
	}

	// ensure load from env vars still succeeds
	if tconf.FromEnv != "42" {
		t.Errorf("got %+v, want %+v", tconf.FromEnv, "42")
	}
}

func TestErrorEnvParse(t *testing.T) {
	os.Setenv("EnvInt", "42.42")
	os.Setenv("EnvFloat", "hello")

	type testconf struct {
		EnvInt   int
		EnvFloat float32
	}
	var tconf testconf
	err := Load(&tconf)

	// check for expected errors
	if err == nil {
		t.Error("got no error, expected parses error")
		return
	}
	if len(err.Fields) != 2 {
		t.Errorf("got %d errors, expected 2 parse errors", len(err.Fields))
		return
	}
	var expectedErr *strconv.NumError
	for _, field := range err.Fields {
		testErr := field.Error

		if !errors.As(testErr, &expectedErr) {
			t.Errorf("got %s, expected strconv.NumError", testErr.Error())
			return
		}
	}

}

func TestErrorNotSupported(t *testing.T) {
	os.Setenv("EnvList", "[42,42]")

	type testconf struct {
		EnvList []int
	}
	var tconf testconf
	err := Load(&tconf)

	// check for expected errors
	if err == nil {
		t.Error("got no error, expected ssconfig.TypeError")
		return
	}
	if len(err.Fields) != 1 {
		t.Errorf("got %d errors, expected 1 error", len(err.Fields))
		return
	}
	var expectedErr *TypeError
	for _, field := range err.Fields {
		testErr := field.Error

		if !errors.As(testErr, &expectedErr) {
			t.Errorf("got '%s', expected ssconfig.TypeError", testErr.Error())
			return
		}
	}
}

func TestErrorPrints(t *testing.T) {
	innerErr := errors.New("Error")
	err := ConfigError{
		[]FieldError{{
			"Field",
			innerErr,
		}},
	}
	expected := "ConfigErrors: [{Field:Field Error:Error}]"
	result := err.Error()

	if result != expected {
		t.Errorf("unexpected error string\ngot:  '%s'\nwant: '%s'", result, expected)
	}

	err2 := TypeError{"TypeName"}
	result = err2.Error()
	expected = "TypeName type not supported"

	if result != expected {
		t.Errorf("unexpected error string\ngot:  '%s'\nwant: '%s'", result, expected)
	}
}
