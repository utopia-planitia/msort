
test:
	diff test-data/golden.yaml <(cat test-data/helmfile.yaml | go run main.go > test-data/golden.yaml)
