# ![datagen logo](assets/cover.png)

`setcfg` (pronounced SET-CONFIG) allows you to set blobs of arbitrary YAML with other blobs of arbitrary YAML. It will walk a YAML file from top to bottom, replacing any blobs of YAML that appear in the file (at any level).

## Installation

`setcfg` can currently be installed via the Go toolchain, although a release will be cut and made available soon:

```
$ go get -u github.com/codingconcepts/setcfg
```

## Usage

**Help text**:
```
setcfg -h
  -i string
        Absolute or relative path to input YAML file.
  -p string
        Absolute or relative path to the parts YAML file.
  -pattern string
        The regex pattern to use for extracting part keys. (default "~(.*?)~")
```

Example:

The following command will place any placeholders found within input.yaml with parts found in parts.yaml:

**input.yaml**:
``` yaml
person:
    name: Rob
    favourite_shows: ~shows~
    pet: ~pet~
```

**parts.yaml**
``` yaml
shows:
- South Park
- Arrested Development

pet:
    name: Twinkle Toes
    age: 2
```

```
$ setcfg -i input.yaml -p parts.yaml

person:
  favourite_shows:
  - South Park
  - Arrested Development
  name: Rob
  pet:
    age: 2
    name: Twinkle Toes
```

`setcfg` outputs to stdout, meaning the results can be piped to a new file or to be included in the results of something like a `kubectl apply` as follows:

**Pipe to file**:
```
$ setcfg -i input.yaml -p parts.yaml > output.yaml
```

**Pipe to kubectl apply**:
```
$ setcfg -i input.yaml -p parts.yaml | kubectl apply -f -
```
