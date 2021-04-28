package runner_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		style          runner.Style
		expectedOutput string
	}{
		{name: "nil style", expectedOutput: "testdata/out.txt"},
		{name: "default style", style: runner.StyleDefault, expectedOutput: "testdata/out.txt"},
		{name: "non-existent style", style: runner.Style("foo"), expectedOutput: "testdata/out.txt"},
		{name: "markdown style", style: runner.StyleMarkdown, expectedOutput: "testdata/out.md"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var gotOut bytes.Buffer

			inBytes, _ := ioutil.ReadFile("testdata/in.txt")
			in := bytes.NewBuffer(inBytes)

			outBytes, _ := ioutil.ReadFile(tt.expectedOutput)
			wantOut := bytes.NewBuffer(outBytes)

			err := runner.Run(in, &gotOut, false, false, false, tt.style)

			if err != nil {
				t.Errorf("Error should be nil, got %w", err)
			}

			if !bytes.Equal(gotOut.Bytes(), wantOut.Bytes()) {
				t.Errorf("Wanted \n%q, got \n%q", wantOut.String(), gotOut.String())
			}
		})
	}
}

func TestRunNoUpdatesCase(t *testing.T) {
	inBytes, err := ioutil.ReadFile("testdata/no_direct_updates.json")
	if err != nil {
		t.Errorf("Failed to read input file: %w", err)
	}
	in := bytes.NewBuffer(inBytes)
	var result bytes.Buffer
	err = runner.Run(in, &result, true, true, false, runner.StyleDefault)
	if err != nil {
		t.Errorf("Error should be nil, got %w", err)
	}
	if result.Len() != 0 {
		t.Errorf("Wanted an empty output, got \n%q", result.String())
	}
}

func TestRunWithError(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/err.txt")
	in := bytes.NewBuffer(inBytes)

	err := runner.Run(in, &out, false, false, false, runner.StyleDefault)

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("Wanted an EOF error, got %w", err)
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

	err := runner.Run(in, &out, false, false, true, runner.StyleDefault)
	if err != nil {
		t.Errorf("Error should be nil, got %w", err)
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

	err := runner.Run(in, &out, false, true, true, runner.StyleDefault)
	if err != nil {
		t.Errorf("Error should be nil, got %s", err.Error())
	}

	if exp := 0; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}
