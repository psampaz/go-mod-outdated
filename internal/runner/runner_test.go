package runner_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		style          runner.OutputStyle
		expectedOutput string
	}{
		{name: "nil style", expectedOutput: "testdata/out.txt"},
		{name: "default style", style: runner.StyleDefault, expectedOutput: "testdata/out.txt"},
		{name: "non-existent style", style: runner.OutputStyle("foo"), expectedOutput: "testdata/out.txt"},
		{name: "markdown style", style: runner.StyleMarkdown, expectedOutput: "testdata/out.md"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var gotOut bytes.Buffer

			inBytes, _ := os.ReadFile("testdata/in.txt")
			in := bytes.NewBuffer(inBytes)

			outBytes, _ := os.ReadFile(tt.expectedOutput)
			wantOut := bytes.NewBuffer(outBytes)

			err := runner.Run(in, &gotOut, false, false, false, tt.style)

			if err != nil {
				t.Errorf("Error should be nil, got %s", err)
			}

			if !bytes.Equal(gotOut.Bytes(), wantOut.Bytes()) {
				t.Errorf("Wanted \n%q, got \n%q", wantOut.String(), gotOut.String())
			}
		})
	}
}

func TestRunNoUpdatesCase(t *testing.T) {
	inBytes, err := os.ReadFile("testdata/no_direct_updates.json")
	if err != nil {
		t.Errorf("Failed to read input file: %s", err)
	}

	in := bytes.NewBuffer(inBytes)

	var out bytes.Buffer

	err = runner.Run(in, &out, true, false, false, runner.StyleDefault)
	if err != nil {
		t.Errorf("Error should be nil, got %s", err)
	}

	if out.Len() != 0 {
		t.Errorf("Wanted an empty output, got \n%q", out.String())
	}
}

func TestRunWithError(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := os.ReadFile("testdata/err.txt")
	in := bytes.NewBuffer(inBytes)

	err := runner.Run(in, &out, false, false, false, runner.StyleDefault)

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("Wanted an EOF error, got %s", err)
	}
}

func TestRunExitWithNonZero(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := os.ReadFile("testdata/in.txt")
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
		t.Errorf("Error should be nil, got %s", err)
	}

	if exp := 1; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}

func TestRunExitWithNonZeroIndirectsOnly(t *testing.T) {
	inBytes, _ := os.ReadFile("testdata/update_indirect.txt")
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
		t.Errorf("Error should be nil, got %s", err)
	}

	if exp := 0; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}
