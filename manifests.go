package main

import (
	"io"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type manifest struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
	yaml string
}

// ByOrder sorts manifests by kind, namespace and name.
type ByOrder []manifest

func (a ByOrder) Len() int      { return len(a) }
func (a ByOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool {
	if a[i].Kind != a[j].Kind {
		return a[i].Kind < a[j].Kind
	}

	if a[i].Metadata.Namespace != a[j].Metadata.Namespace {
		return a[i].Metadata.Namespace < a[j].Metadata.Namespace
	}

	if a[i].Metadata.Name != a[j].Metadata.Name {
		return a[i].Metadata.Name < a[j].Metadata.Name
	}

	return a[i].yaml < a[j].yaml
}

func splitManifests(yml string) []manifest {
	yml = strings.TrimSpace(yml)

	if strings.HasPrefix(yml, "---\n") {
		yml = "\n" + yml
	}

	if strings.HasSuffix(yml, "\n---") {
		yml += "\n"
	}

	separator := regexp.MustCompile("\n---\n")
	chunks := separator.Split(yml, -1)
	manifests := []manifest{}

	for _, doc := range chunks {
		doc = strings.TrimSpace(doc)

		if doc == "" {
			continue
		}

		manifest, err := parseManifest(doc)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Panicln(err)
		}

		manifests = append(manifests, manifest)
	}

	return manifests
}

func parseManifest(yml string) (manifest, error) {
	d := yaml.NewDecoder(strings.NewReader(yml))
	m := manifest{}

	err := d.Decode(&m) // invalid yaml returns empty string as a kind
	if err != nil {
		return manifest{}, err
	}

	m.yaml = yml

	return m, nil
}
