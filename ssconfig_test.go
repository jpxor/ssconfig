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
	"os"
	"reflect"
	"testing"
)

func TestSuccessPaths(t *testing.T) {
	type testConfig struct {

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

	var tconf testConfig
	Set{FilePath: "./test.json"}.Load(&tconf)

	answerInt := 42
	answerFloat := 4.2
	answerString := "42"
	answerList := []float32{4.2, 4.2, 4.2}
	answerMap := map[string]int{"answer": 42}

	t.Run("Load from File", func(t *testing.T) {
		if tconf.FileInt != answerInt {
			t.Errorf("got %d, want %d", tconf.FileInt, answerInt)
		}
		if tconf.FileString != answerString {
			t.Errorf("got '%s', want '%s'", tconf.FileString, answerString)
		}
		if !reflect.DeepEqual(tconf.FileList, answerList) {
			t.Errorf("got '%+v', want '%+v'", tconf.FileList, answerList)
		}
		if !reflect.DeepEqual(tconf.FileMap, answerMap) {
			t.Errorf("got '%+v', want '%+v'", tconf.FileMap, answerMap)
		}
	})

	t.Run("Load from Env", func(t *testing.T) {
		if tconf.EnvVarInt != answerInt {
			t.Errorf("got %d, want %d", tconf.EnvVarInt, answerInt)
		}
		if tconf.EnvVarInt8 != int8(answerInt) {
			t.Errorf("got %d, want %d", tconf.EnvVarInt8, answerInt)
		}
		if tconf.EnvVarInt16 != int16(answerInt) {
			t.Errorf("got %d, want %d", tconf.EnvVarInt16, answerInt)
		}
		if tconf.EnvVarInt32 != int32(answerInt) {
			t.Errorf("got %d, want %d", tconf.EnvVarInt32, answerInt)
		}
		if tconf.EnvVarInt64 != int64(answerInt) {
			t.Errorf("got %d, want %d", tconf.EnvVarInt64, answerInt)
		}

		if tconf.EnvVarString != answerString {
			t.Errorf("got '%s', want '%s'", tconf.EnvVarString, answerString)
		}

		if tconf.EnvVarFloat32 != float32(answerFloat) {
			t.Errorf("got %f, want %f", tconf.EnvVarFloat32, answerFloat)
		}
		if tconf.EnvVarFloat64 != answerFloat {
			t.Errorf("got %f, want %f", tconf.EnvVarFloat64, answerFloat)
		}
	})

	t.Run("Env Overwrites File", func(t *testing.T) {
		if tconf.OverwriteInt != answerInt {
			t.Errorf("got %d, want %d", tconf.EnvVarInt, answerInt)
		}
		if tconf.OverwriteInt8 != int8(answerInt) {
			t.Errorf("got %d, want %d", tconf.OverwriteInt8, answerInt)
		}
		if tconf.OverwriteInt16 != int16(answerInt) {
			t.Errorf("got %d, want %d", tconf.OverwriteInt16, answerInt)
		}
		if tconf.OverwriteInt32 != int32(answerInt) {
			t.Errorf("got %d, want %d", tconf.OverwriteInt32, answerInt)
		}
		if tconf.OverwriteInt64 != int64(answerInt) {
			t.Errorf("got %d, want %d", tconf.OverwriteInt64, answerInt)
		}

		if tconf.OverwriteString != answerString {
			t.Errorf("got '%s', want '%s'", tconf.OverwriteString, answerString)
		}

		if tconf.OverwriteFloat32 != float32(answerFloat) {
			t.Errorf("got %f, want %f", tconf.OverwriteFloat32, answerFloat)
		}
		if tconf.OverwriteFloat64 != answerFloat {
			t.Errorf("got %f, want %f", tconf.OverwriteFloat64, answerFloat)
		}
	})
}
