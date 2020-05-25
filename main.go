package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	yml, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err.Error())
	}

	manifests := splitManifests(string(yml))

	sort.Sort(ByOrder(manifests))

	for _, manifest := range manifests {
		_, err := fmt.Printf("---\n%s\n", manifest.yaml)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
