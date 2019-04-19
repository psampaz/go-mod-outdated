package runner

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

func TestRun(t *testing.T) {

	var gotOut bytes.Buffer

	inBytes, _ := ioutil.ReadFile("testdata/in.txt")
	in := bytes.NewBuffer(inBytes)

	outBytes, _ := ioutil.ReadFile("testdata/out.txt")
	wantOut := bytes.NewBuffer(outBytes)

	err := Run(in, &gotOut, false, false)

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

	gotErr := Run(in, &out, false, false)
	wantErr := errors.New("unexpected EOF")

	if gotErr.Error() != wantErr.Error() {
		t.Errorf("Wanted %s, got %s", wantErr, gotErr)
	}

}
