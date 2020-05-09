# msort

msort takes a list of kubernetes manifests as input and sorts them to produce a reproducible artifact.

## Example

``` bash
helmfile template | msort > manifests.yaml
```
