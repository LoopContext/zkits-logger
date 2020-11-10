// Copyright 2020 The ZKits Project Authors.
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	o := New("test")

	if o == nil {
		t.Fatal("New(): nil")
	}
}

func TestLogger_Name(t *testing.T) {
	o := New("test")

	if name := o.Name(); name != "test" {
		t.Fatalf("Logger.Name(): %s", name)
	}
}

func TestLogger_Level(t *testing.T) {
	o := New("test")

	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
	if o.SetLevel(DebugLevel) == nil {
		t.Fatal("Logger.SetLevel() return nil")
	}
	if level := o.GetLevel(); level != DebugLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
}

func TestLogger_Output(t *testing.T) {
	o := New("test")
	w := new(bytes.Buffer)

	if o.SetOutput(w) == nil {
		t.Fatal("Logger.SetOutput(): nil")
	}
	if o.SetOutput(nil) == nil {
		t.Fatal("Logger.SetOutput(nil): nil")
	}
}

func TestLogger_SetNowFunc(t *testing.T) {
	o := New("test")
	f := func() time.Time { return time.Now() }

	if o.SetNowFunc(f) == nil {
		t.Fatal("Logger.SetNowFunc(): nil")
	}
	if o.SetNowFunc(nil) == nil {
		t.Fatal("Logger.SetNowFunc(nil): nil")
	}
}

func TestLogger_SetExitFunc(t *testing.T) {
	o := New("test")
	f := func(int) {}

	if o.SetExitFunc(f) == nil {
		t.Fatal("Logger.SetExitFunc(): nil")
	}
	if o.SetExitFunc(nil) == nil {
		t.Fatal("Logger.SetExitFunc(nil): nil")
	}
}

func TestLogger_SetPanicFunc(t *testing.T) {
	o := New("test")
	f := func(string) {}

	if o.SetPanicFunc(f) == nil {
		t.Fatal("Logger.SetPanicFunc(): nil")
	}
	if o.SetPanicFunc(nil) == nil {
		t.Fatal("Logger.SetPanicFunc(nil): nil")
	}
}

func TestLogger_SetFormatter(t *testing.T) {
	o := New("test")
	f := FormatterFunc(func(Entity, *bytes.Buffer) error { return nil })

	if o.SetFormatter(f) == nil {
		t.Fatal("Logger.SetFormatter(): nil")
	}
	if o.SetFormatter(nil) == nil {
		t.Fatal("Logger.SetFormatter(nil): nil")
	}
}

func TestLogger_Log(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")

	o.SetOutput(w)
	o.SetExitFunc(nil)  // Disable
	o.SetPanicFunc(nil) // Disable
	o.SetLevel(TraceLevel)

	now := time.Now()
	o.SetNowFunc(func() time.Time { return now })

	do := func(s string, fs ...func(Logger) (Level, string)) {
		buf := new(bytes.Buffer)
		for _, f := range fs {
			w.Reset()
			buf.Reset()
			level, message := f(o)
			if level.IsValid() {
				err := json.NewEncoder(buf).Encode(map[string]interface{}{
					"name":    "test",
					"time":    now.Format(time.RFC3339),
					"level":   level.String(),
					"message": message,
					"fields":  map[string]interface{}{},
				})
				if err != nil {
					t.Fatalf("%s: %s", s, err)
				}
				if !bytes.Equal(w.Bytes(), buf.Bytes()) {
					t.Fatalf("%s: %s -- %s", s, w.String(), buf.String())
				}
			} else {
				// No log
				if got := w.String(); got != "" {
					t.Fatalf("%s: %s", s, got)
				}
			}
		}
	}

	// -------------- TraceLevel -----------------

	do("TraceLevel", func(o Logger) (Level, string) {
		o.Trace("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Traceln("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Tracef("foo-%s", "bar")
		return TraceLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Print("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Println("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Printf("foo-%s", "bar")
		return TraceLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- DebugLevel -----------------

	do("DebugLevel", func(o Logger) (Level, string) {
		o.Debug("foo")
		return DebugLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Debugln("foo")
		return DebugLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Debugf("foo-%s", "bar")
		return DebugLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- InfoLevel -----------------

	do("InfoLevel", func(o Logger) (Level, string) {
		o.Info("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Infoln("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Infof("foo-%s", "bar")
		return InfoLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Echo("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Echoln("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Echof("foo-%s", "bar")
		return InfoLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- WarnLevel -----------------

	do("WarnLevel", func(o Logger) (Level, string) {
		o.Warn("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warnln("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warnf("foo-%s", "bar")
		return WarnLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Warning("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warningln("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warningf("foo-%s", "bar")
		return WarnLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- ErrorLevel -----------------

	do("ErrorLevel", func(o Logger) (Level, string) {
		o.Error("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Errorln("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Errorf("foo-%s", "bar")
		return ErrorLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- FatalLevel -----------------

	do("FatalLevel", func(o Logger) (Level, string) {
		o.Fatal("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatalln("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatalf("foo-%s", "bar")
		return FatalLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- PanicLevel -----------------

	do("PanicLevel", func(o Logger) (Level, string) {
		o.Panic("foo")
		return PanicLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panicln("foo")
		return PanicLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panicf("foo-%s", "bar")
		return PanicLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// Test a higher level of log.
	o.SetLevel(ErrorLevel)

	do("Use ErrorLevel", func(o Logger) (Level, string) {
		o.Trace("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Debug("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Info("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Warn("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Error("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatal("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panic("foo")
		return PanicLevel, "foo"
	})
}

func TestLogger_With(t *testing.T) {

}

func TestLogger_Hook(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var message string

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		message = s.Message()
		return nil
	})

	o.Trace("foo")

	if message != "foo" {
		t.Fatalf("Log hook not exec: %s", message)
	}
}

type testErrorJSONMarshaler string

func (s testErrorJSONMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New(string(s))
}

type testErrorWriter string

func (s testErrorWriter) Write([]byte) (int, error) {
	return 0, errors.New(string(s))
}
