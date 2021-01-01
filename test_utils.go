package tgin

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expect, value interface{}, msgs ...string) {
	if !reflect.DeepEqual(expect, value) {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect: %+v but got: %+v", msg, expect, value)
	}
}

func assertNotEqual(t *testing.T, expect, value interface{}, msgs ...string) {
	if reflect.DeepEqual(expect, value) {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect: %+v but got: %+v", msg, expect, value)
	}
}

func assertTrue(t *testing.T, value bool, msgs ...string) {
	if !value {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%svalue is %v", msg, value)
	}
}

func assertFalse(t *testing.T, value bool, msgs ...string) {
	assertTrue(t, !value, msgs...)
}

func assertNil(t *testing.T, value interface{}, msgs ...string) {
	if value != nil {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect nil but got: %+v", msg, value)
	}
}

func assertNotNil(t *testing.T, value interface{}, msgs ...string) {
	if value == nil {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect not nil but got nil", msg)
	}
}

func assertBody(t *testing.T, resp *http.Response, expect string, msgs ...string) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Read body got error: %v", err)
	}

	if string(data) != expect {
		msg := strings.TrimSpace(strings.Join(msgs, " "))
		if msg != "" {
			msg += "; "
		}
		t.Fatalf("%sExpect %s but got %s", msg, expect, string(data))
	}
}
