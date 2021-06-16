package runner_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/psampaz/go-mod-outdated/internal/mod"
	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		style          string
		expectedOutput string
	}{
		{name: "nil style", expectedOutput: "testdata/out.txt"},
		{name: "default style", style: "default", expectedOutput: "testdata/out.txt"},
		{name: "non-existent style", style: "foo", expectedOutput: "testdata/out.txt"},
		{name: "markdown style", style: "markdown", expectedOutput: "testdata/out.md"},
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

	gotErr := runner.Run(in, &out, false, false, false, "default")
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

	err := runner.Run(in, &out, false, false, true, "default")
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

	err := runner.Run(in, &out, false, true, true, "default")
	if err != nil {
		t.Errorf("Error should be nil, got %s", err.Error())
	}

	if exp := 0; got != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, got)
	}
}

func TestHTMLTable(t *testing.T) {
	var actualOutput bytes.Buffer
	moduleInput := []mod.Module{mod.Module{
		Path:    "github.com/mattn/go-runewidth",
		Version: "v0.0.10",
		Update: &mod.Module{
			Version: "v0.0.12",
		},
		Indirect: true,
	}}
	err := runner.RenderHTMLTable(&actualOutput, moduleInput)
	if err != nil {
		t.Errorf("Error should be nil, got %w", err)
	}
	expectedBytes, err := ioutil.ReadFile("testdata/expected_table.html")
	expectedOutput := bytes.NewBuffer(expectedBytes)

	if actualOutput.String() != expectedOutput.String() {
		t.Errorf("Expected table output to match \n%v, but got \n%v", expectedOutput, actualOutput.String())
	}
}
