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
		name    string
		in      string
		golden  string
		keySort bool
	}{
		{
			name:    "sort output of helmfile template",
			in:      "testdata/case1.yaml",
			golden:  "testdata/case1.golden.yaml",
			keySort: false,
		},
		{
			name:    "sort keys for none kubernetes manifests",
			in:      "testdata/case2.yaml",
			golden:  "testdata/case2.golden.yaml",
			keySort: true,
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

			os.Args = []string{rescueArgs[0], tt.in}

			os.Setenv("DISABLE_KEY_SORTING", "1")
			if tt.keySort {
				os.Setenv("DISABLE_KEY_SORTING", "")
			}

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
