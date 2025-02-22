// Copyright 2021 The ZKits Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewFileWriter(t *testing.T) {
	dir := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "TestNewFileWriter")
	if err := os.MkdirAll(dir, 0766); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	name := filepath.Join(dir, "test.log")

	if w, err := NewFileWriter(name, 1024, 0); err != nil {
		t.Fatal(err)
	} else {
		if w == nil {
			t.Fatal("NewFileWriter(): nil")
		}
		if err = w.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := NewFileWriter(dir, 1024, 0); err == nil {
		t.Fatal("NewFileWriter(): nil error")
	}
}

func TestMustNewFileWriter(t *testing.T) {
	dir := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "TestMustNewFileWriter")
	if err := os.MkdirAll(dir, 0766); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	name := filepath.Join(dir, "test.log")
	if w := MustNewFileWriter(name, 1024, 0); w == nil {
		t.Fatal("MustNewFileWriter(): nil")
	} else {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}

	defer func() {
		if v := recover(); v == nil {
			t.Fatal("MustNewFileWriter(): not panic")
		}
	}()
	MustNewFileWriter(dir, 1024, 0)
}

func TestFileWriter(t *testing.T) {
	dir := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "TestFileWriter")
	if err := os.MkdirAll(dir, 0766); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	name := filepath.Join(dir, "test.log")
	w := MustNewFileWriter(name, 1024, 0)
	defer func() {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	data := strings.Repeat("1", 1024)
	if n, err := w.Write([]byte(data)); err != nil {
		t.Fatal(err)
	} else {
		if n != 1024 {
			t.Fatalf("FileWriter.Write(): %d", n)
		}
	}
	if matches, err := filepath.Glob(name + ".*"); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) == 0 {
			t.Fatal("FileWriter.Write(): not rename")
		}
		got, err := ioutil.ReadFile(matches[0])
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != data {
			t.Fatalf("FileWriter.Write(): got %s", string(got))
		}
	}
}

func TestFileWriterWithBackup(t *testing.T) {
	dir := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10), "TestFileWriterWithBackup")
	if err := os.MkdirAll(dir, 0766); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	name := filepath.Join(dir, "test.log")
	w := MustNewFileWriter(name, 1000, 2)
	defer func() {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	data := strings.Repeat("1", 1024)
	for i := 0; i < 3; i++ {
		if n, err := w.Write([]byte(data)); err != nil {
			t.Fatal(err)
		} else {
			if n != 1024 {
				t.Fatalf("FileWriter.Write(): %d", n)
			}
		}
	}
	// The file system takes some time.
	time.Sleep(time.Second)
	if matches, err := filepath.Glob(name + ".*"); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) == 0 {
			t.Fatal("FileWriter.Write(): not rename")
		}
		if len(matches) != 2 {
			t.Fatalf("FileWriter.Write(): not clean up: %v", matches)
		}
	}
}
