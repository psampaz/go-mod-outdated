package runner_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		markdown bool
		expected string
	}{
		{name: "ascii table", markdown: false, expected: "testdata/out.txt"},
		{name: "markdown table", markdown: true, expected: "testdata/out.md"},
	}
	//scopelint:ignore
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotOut bytes.Buffer

			inBytes, _ := ioutil.ReadFile("testdata/in.txt")
			in := bytes.NewBuffer(inBytes)

			outBytes, _ := ioutil.ReadFile(tt.expected)
			wantOut := bytes.NewBuffer(outBytes)

			err := runner.Run(in, &gotOut, false, false, false, tt.markdown)

			if err != nil {
				t.Errorf("Error should be nil, got %s", err.Error())
			}

			if !bytes.Equal(gotOut.Bytes(), wantOut.Bytes()) {
				t.Errorf("Wanted \n%q, got \n%q", wantOut.String(), gotOut.String())
			}
		})
	}
}

func TestRunWithError(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/err.txt")
	in := bytes.NewBuffer(inBytes)

	gotErr := runner.Run(in, &out, false, false, false, false)
	wantErr := errors.New("unexpected EOF")

	if gotErr.Error() != wantErr.Error() {
		t.Errorf("Wanted %q, got %q", wantErr, gotErr)
	}
}

func TestRunExitWithNonZero(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/in.txt")
	in := bytes.NewBuffer(inBytes)

	oldOsExit := runner.OsExit

	defer func() { runner.OsExit = oldOsExit }()

	var got int

	testExit := func(code int) {
		got = code
	}

	runner.OsExit = testExit

	err := runner.Run(in, &out, false, false, true, false)
	if err != nil {
		t.Errorf("Error should be nil, got %s", err.Error())
	}

	if exp := 1; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}

func TestRunExitWithNonZeroIndirectsOnly(t *testing.T) {
	inBytes, _ := ioutil.ReadFile("testdata/update_indirect.txt")
	in := bytes.NewBuffer(inBytes)

	oldOsExit := runner.OsExit

	defer func() { runner.OsExit = oldOsExit }()

	var got int

	testExit := func(code int) {
		got = code
	}

	runner.OsExit = testExit

	var out bytes.Buffer

	err := runner.Run(in, &out, false, true, true, false)
	if err != nil {
		t.Errorf("Error should be nil, got %s", err.Error())
	}

	if exp := 0; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}
