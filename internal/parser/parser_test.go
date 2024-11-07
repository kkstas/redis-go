package parser_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kkstas/redis-go/internal/parser"
)

func TestParseArray(t *testing.T) {
	cases := []struct {
		input   string
		want    []string
		isValid bool
	}{
		{"*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$11\r\nhello world\r\n$2\r\npx\r\n$4\r\n1000\r\n", []string{"SET", "key", "hello world", "px", "1000"}, true},
		{"*2\r\n$4\r\nECHO\r\n$5\r\ngrape\r\n", []string{"ECHO", "grape"}, true},
		{"*2\r\n$4\r\nECHO\r\n$14\r\nasdf asdf\r\n", []string{"ECHO", "asdf asdf"}, true},
		{"*1\r\n$4\r\nasdf\r\n", []string{"asdf"}, true},
		{"*3\r\n$3\r\nSET\r\n$4\r\naaa\r\n$3\r\nxyz\r\n", []string{"SET", "aaa", "xyz"}, true},
		{"*2\r\n$3\r\nGET\r\n$3\r\naaa\r\n", []string{"GET", "aaa"}, true},
		{"*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\npx\r\n$3\r\n100\r\n", []string{"SET", "foo", "bar", "px", "100"}, true},
		{"*\r\nasdfasdf\r\nasdfasdf\r\n", nil, false},
		{"xd\r\nasdf\r\nasdf\r\n", nil, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("parses string to %v", c.want), func(t *testing.T) {
			got, err := parser.Parse(c.input)

			if c.isValid && err != nil {
				t.Errorf("didn't expect an error but got one: %v", err)
			}
			if !c.isValid && err == nil {
				t.Error("expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %#v want %#v", got, c.want)
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	t.Run("recognizes PING case-insensitive", func(t *testing.T) {
		want := "ping"
		assertEqual(t, parser.GetCommand([]string{"PING"}), want)
		assertEqual(t, parser.GetCommand([]string{"piNG"}), want)
		assertEqual(t, parser.GetCommand([]string{"ping"}), want)
	})

	t.Run("recognizes ECHO case-insensitive", func(t *testing.T) {
		want := "echo"
		assertEqual(t, parser.GetCommand([]string{"ECHO", "grape grape"}), want)
		assertEqual(t, parser.GetCommand([]string{"Echo", "grape grape"}), want)
		assertEqual(t, parser.GetCommand([]string{"echo", "grape grape"}), want)
	})

	t.Run("doesn't recognize ECHO if arr is too long/too short", func(t *testing.T) {
		want := ""
		assertEqual(t, parser.GetCommand([]string{"ECHO"}), want)
		assertEqual(t, parser.GetCommand([]string{"ECHO", "asdf", "asdf"}), want)
	})

	t.Run("recognizes SET case-insensitive", func(t *testing.T) {
		want := "set"
		assertEqual(t, parser.GetCommand([]string{"SET", "aaa", "xyz"}), want)
		assertEqual(t, parser.GetCommand([]string{"sET", "aaa", "xyz"}), want)
		assertEqual(t, parser.GetCommand([]string{"set", "aaa", "xyz"}), want)
	})

	t.Run("doesn't recognize SET if arr is too long/too short", func(t *testing.T) {
		want := ""
		assertEqual(t, parser.GetCommand([]string{"SET", "aaa", "xyz", "aaa", "xyz"}), want)
		assertEqual(t, parser.GetCommand([]string{"SET", "aaa", "xyz", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"SET", "aaa"}), want)
	})

	t.Run("recognizes SET with expiry case-insensitive", func(t *testing.T) {
		want := "set_expiry"
		assertEqual(t, parser.GetCommand([]string{"SET", "foo", "bar", "px", "100"}), want)
		assertEqual(t, parser.GetCommand([]string{"seT", "foo", "bar", "px", "100"}), want)
		assertEqual(t, parser.GetCommand([]string{"set", "foo", "bar", "px", "100"}), want)
	})

	t.Run("doesn't recognize SET with expiry if px is not 4th element", func(t *testing.T) {
		want := ""
		assertEqual(t, parser.GetCommand([]string{"set", "foo", "bar", "xx", "100"}), want)
		assertEqual(t, parser.GetCommand([]string{"SET", "aaa", "xyz", "", "xyz"}), want)
	})

	t.Run("recognizes GET case-insensitive", func(t *testing.T) {
		want := "get"
		assertEqual(t, parser.GetCommand([]string{"GET", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"gET", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"get", "aaa"}), want)
	})

	t.Run("doesn't recognize GET if arr is too long/too short", func(t *testing.T) {
		want := ""
		assertEqual(t, parser.GetCommand([]string{"GET", "aaa", "xyz", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"GET"}), want)
	})

	t.Run("recognizes CONFIG GET case-insensitive", func(t *testing.T) {
		want := "config_get"
		assertEqual(t, parser.GetCommand([]string{"CONFIG", "GET", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"ConfIg", "geT", "aaa"}), want)
		assertEqual(t, parser.GetCommand([]string{"config", "get", "aaa"}), want)
	})

	t.Run("doesn't recognize CONFIG GET if arr is too long/too short", func(t *testing.T) {
		want := ""
		assertEqual(t, parser.GetCommand([]string{"CONFIG", "GET", "aaa", "sdfa"}), want)
		assertEqual(t, parser.GetCommand([]string{"CONFIG", "GET"}), want)
	})
}

func TestToSimpleString(t *testing.T) {
	got := parser.ToSimpleString("some message")
	want := "+some message\r\n"
	assertEqual(t, got, want)

	got = parser.ToSimpleString("")
	want = "+\r\n"
	assertEqual(t, got, want)
}

func TestToBulkString(t *testing.T) {
	got := parser.ToBulkString("some message")
	want := "$12\r\nsome message\r\n"
	assertEqual(t, got, want)

	got = parser.ToBulkString("")
	want = "$0\r\n\r\n"
	assertEqual(t, got, want)
}

func TestToRESPArray(t *testing.T) {
	t.Run("assembles RESP array", func(t *testing.T) {
		got := parser.ToRESPArray([]string{"dir", "/tmp/redis-files"})
		want := "*2\r\n$3\r\ndir\r\n$16\r\n/tmp/redis-files\r\n"
		if got != want {
			t.Errorf("got %s, want %s", parser.GetRaw(got), parser.GetRaw(want))
		}

	})
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't expect one: %v", err)
	}
}

func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
