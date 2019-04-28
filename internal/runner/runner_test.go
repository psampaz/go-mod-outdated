package runner_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func TestRun(t *testing.T) {

	var gotOut bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/in.txt")
	in := bytes.NewBuffer(inBytes)

	outBytes, _ := ioutil.ReadFile("testdata/out.txt")
	wantOut := bytes.NewBuffer(outBytes)

	err := runner.Run(in, &gotOut, false, false, false)

	if err != nil {
		t.Errorf("Error should be nil, got %s", err.Error())
	}

	if !bytes.Equal(gotOut.Bytes(), wantOut.Bytes()) {
		t.Errorf("Wanted \n%s, got \n%s", wantOut.String(), gotOut.String())
	}

}

func TestRunWithError(t *testing.T) {

	var out bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/err.txt")
	in := bytes.NewBuffer(inBytes)

	gotErr := runner.Run(in, &out, false, false, false)
	wantErr := errors.New("unexpected EOF")

	if gotErr.Error() != wantErr.Error() {
		t.Errorf("Wanted %s, got %s", wantErr, gotErr)
	}

}


func TestRunExitWithNonZero(t *testing.T) {
	var out bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/in.txt")
	in := bytes.NewBuffer(inBytes)

	if os.Getenv("TEST_EXITCODE") == "1" {
		runner.Run(in, &out, false, false, true)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRunExitWithNonZero")
	cmd.Env = append(os.Environ(), "TEST_EXITCODE=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}