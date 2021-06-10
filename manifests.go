package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
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
			log.Panicf("parse yaml failed: %v: %v", err, doc)
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

func (m manifest) Print() (string, error) {
	return m.yaml, nil
}

func SortedByKeys(in string) (string, error) {
	dec := yaml.NewDecoder(strings.NewReader(in))

	var doc yaml.Node
	err := dec.Decode(&doc)
	if err != nil {
		return "", fmt.Errorf("decode yaml: %v", err)
	}

	sortYAML(&doc)

	out := bytes.NewBuffer(nil)

	enc := yaml.NewEncoder(out)
	defer enc.Close()

	enc.SetIndent(2)

	err = enc.Encode(&doc)
	if err != nil {
		return "", fmt.Errorf("encode yaml: %v", err)
	}

	return out.String(), nil
}

type byKey []*yaml.Node

func (i byKey) Len() int { return len(i) / 2 }

func (i byKey) Swap(x, y int) {
	x *= 2
	y *= 2
	i[x], i[y] = i[y], i[x]         // keys
	i[x+1], i[y+1] = i[y+1], i[x+1] // values
}

func (i byKey) Less(x, y int) bool {
	x *= 2
	y *= 2
	return i[x].Value < i[y].Value
}

func sortYAML(node *yaml.Node) *yaml.Node {
	for i, n := range node.Content {
		node.Content[i] = sortYAML(n)
	}

	if node.Kind == yaml.MappingNode {
		sort.Sort(byKey(node.Content))
	}

	return node
}
