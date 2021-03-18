# msort

msort takes a list of kubernetes manifests as input and sorts them to produce a reproducible artifact.

## sorted keys

To generate always the same output, the manifests are sorted by key.
The sorting happens by using [yq](https://github.com/mikefarah/yq). yq is a wrapper around [jq](https://github.com/stedolan/jq).

To disable this feature set `DISABLE_KEY_SORTING=1`.

## Example

``` bash
helmfile template | msort > manifests.yaml
```
