package main

import (
	"bytes"
	_ "embed"
	"io/ioutil"
	"os"
	"testing"

	"github.com/andreyvit/diff"
)

//go:embed testdata/case1.golden.yaml
var goldenCase1 []byte

func Test_main(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		golden    string
		keySort   bool
		dropTests bool
	}{
		{
			name:   "sort output of helmfile template",
			in:     "testdata/case1.yaml",
			golden: "testdata/case1.golden.yaml",
		},
		{
			name:    "sort keys for none kubernetes manifests",
			in:      "testdata/case2.yaml",
			golden:  "testdata/case2.golden.yaml",
			keySort: true,
		},
		{
			name:      "drop only test",
			in:        "testdata/case3.yaml",
			golden:    "testdata/case3.golden.yaml",
			dropTests: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			rescueArgs := os.Args
			defer func() {
				os.Stdout = rescueStdout
				os.Args = rescueArgs
			}()

			r, w, _ := os.Pipe()
			os.Stdout = w

			os.Args = []string{rescueArgs[0]}

			if tt.keySort {
				os.Args = append(os.Args, "--sort-keys")
			}

			if tt.dropTests {
				os.Args = append(os.Args, "--drop-tests")
			}

			os.Args = append(os.Args, tt.in)

			main()

			w.Close()
			out, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatalf("read generated output: %v", err)
			}

			golden, err := ioutil.ReadFile(tt.golden)
			if err != nil {
				t.Fatalf("read golden file %s: %v", tt.golden, err)
			}

			if !bytes.Equal(out, golden) {
				t.Logf("generated output of %s does not match content of %s:\n%v", tt.in, tt.golden, diff.LineDiff(string(golden), string(out)))
				t.Fail()
			}
		})
	}
}
