package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	msort "github.com/utopia-planitia/msort/pkg"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run() error {
	yml := []byte{}
	var err error = nil

	if len(os.Args) == 1 {
		yml, err = ioutil.ReadAll(os.Stdin)
	} else {
		yml, err = ioutil.ReadFile(os.Args[1])
	}
	if err != nil {
		return err
	}

	manifests := msort.SplitManifests(string(yml))

	if os.Getenv("DISABLE_KEY_SORTING") == "" {
		for i := range manifests {
			manifests[i].Yaml, err = msort.SortedByKeys(manifests[i].Yaml)
			if err != nil {
				return err
			}
		}
	}

	sort.Sort(msort.ByOrder(manifests))

	first := true

	for _, manifest := range manifests {
		if os.Getenv("KEEP_TESTS") == "" {
			if strings.Contains(strings.ToLower(manifest.Metadata.Name), "test") {
				continue
			}
		}

		yaml, err := manifest.Print()
		if err != nil {
			return err
		}

		if strings.TrimSpace(yaml) == "" {
			continue
		}

		if !first {
			_, err = fmt.Print("---\n")
			if err != nil {
				return err
			}
		}

		first = false

		_, err = fmt.Printf("%s\n", strings.TrimSpace(yaml))
		if err != nil {
			return err
		}
	}

	return nil
}
