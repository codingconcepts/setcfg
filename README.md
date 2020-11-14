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
region: ~region~

credentials:
    username: ~username~
    password: ~password~

brokers:
    - broker: ~broker-1~
    - broker: ~broker-2~

subnet_cidrs: ~subnet-cidrs~
```

**parts.yaml**
``` yaml
region: eu-west-1

username: admin
password: supersecret

broker-1:
    host: https://localhost
    port: 8001

broker-2:
    host: https://localhost
    port: 8002

subnet-cidrs:
- 1.2.3.0/25
- 1.2.3.128/25
```

```
$ setcfg -i input.yaml -p parts.yaml

brokers:
- broker:
    host: https://localhost
    port: 8001
- broker:
    host: https://localhost
    port: 8002
credentials:
  password: supersecret
  username: admin
region: eu-west-1
subnet_cidrs:
- 1.2.3.0/25
- 1.2.3.128/25
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
