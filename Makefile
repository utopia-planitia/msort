
test:
	bash -c "diff test-data/golden.yaml <(cat test-data/helmfile.yaml | DISABLE_KEY_SORTING=1 go run .)"
