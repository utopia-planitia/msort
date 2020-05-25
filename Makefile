
test:
	bash -c "diff test-data/golden.yaml <(cat test-data/helmfile.yaml | go run . > test-data/golden.yaml)"
