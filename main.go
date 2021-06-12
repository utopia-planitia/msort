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
	yml := []byte{}
	var err error = nil

	if len(os.Args) == 1 {
		yml, err = ioutil.ReadAll(os.Stdin)
	} else {
		yml, err = ioutil.ReadFile(os.Args[1])
	}
	if err != nil {
		log.Fatalln(err.Error())
	}

	manifests := msort.SplitManifests(string(yml))

	if os.Getenv("DISABLE_KEY_SORTING") == "" {
		for i := range manifests {
			manifests[i].Yaml, err = msort.SortedByKeys(manifests[i].Yaml)
			if err != nil {
				log.Fatalln(err.Error())
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
			log.Fatalln(err.Error())
		}

		if strings.TrimSpace(yaml) == "" {
			continue
		}

		if !first {
			_, err = fmt.Print("---\n")
			if err != nil {
				log.Fatalln(err.Error())
			}
		}

		first = false

		_, err = fmt.Printf("%s\n", strings.TrimSpace(yaml))
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
