package tgin

import (
	"reflect"
	"strings"
	"testing"
)

func AssertEqual(t *testing.T, expect, value interface{}, msgs ...string) {
	if !reflect.DeepEqual(expect, value) {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect: %+v but got: %+v", msg, expect, value)
	}
}

func AssertNotEqual(t *testing.T, expect, value interface{}, msgs ...string) {
	if reflect.DeepEqual(expect, value) {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect: %+v but got: %+v", msg, expect, value)
	}
}

func AssertTrue(t *testing.T, value bool, msgs ...string) {
	if value {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%svalue is %v", msg, value)
	}
}

func AssertFalse(t *testing.T, value bool, msgs ...string) {
	AssertTrue(t, !value, msgs...)
}

func AssertNil(t *testing.T, value interface{}, msgs ...string) {
	if value != nil {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect nil but got: %+v", msg, value)
	}
}

func AssertNotNil(t *testing.T, value interface{}, msgs ...string) {
	if value == nil {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect not nil but got nil", msg)
	}
}
