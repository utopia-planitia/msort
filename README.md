# msort

msort takes a list of kubernetes manifests as input and sorts them to produce a reproducible artifact.

## Help

``` bash
NAME:
   msort - sort yaml manifests

USAGE:
   msort [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --sort-keys     sort keys within each yaml document (default: false)
   --drop-tests    remove yaml documents with "test" in its name (default: false)
   --in-place, -i  update files in place (default: false)
   --help, -h      show help (default: false)

```

## Example

Sort yaml documents within `manifests.yaml` by Kind, Name, and Namspace, then print the result to `stdout`.

``` bash
msort --in-place --sort-keys file.yaml
```

Read `file.yaml`, sort by Kind, Name, and Namespace, then and sort maps by key and update the file.

``` bash
msort --in-place --sort-keys file.yaml
```

Sort the output of `helmfile template` by Kind, Name, and Namespace, then remove kubernetes manifests with `test` in their name, then write to `stdout`.

``` bash
helmfile template | msort --drop-tests > manifests.yaml
```
