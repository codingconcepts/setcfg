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
  -e string
        Absolute or relative path to the environment YAML file.
  -f value
        A list of 'key=value' fields to substitute (useful as an alternative to -e if all you're substituting are simple fields).
  -i string
        Absolute or relative path to input YAML file.
  -pattern string
        The regex pattern to use for extracting keys. (default "~(.*?)~")
```

Example:

The following command will place any placeholders found within input.yaml with fields found in dev.yaml:

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

**dev.yaml**
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
$ setcfg -i input.yaml -e dev.yaml

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

You can set ad-hoc fields to add to, or override any fields in the env file:
```
$ setcfg -i input.yaml -e dev.yaml -f region=eu-west-2

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
region: eu-west-2
subnet_cidrs:
- 1.2.3.0/25
- 1.2.3.128/25
```

## Multi-document files

`setcfg` supports multi-document files by default. For example:

input.yaml:
``` yaml
a: ~a~
---
b: ~b~
---
c: ~c~
```

```
$ setcfg -i input.yaml -f a=1 -f b=2 -f c=3
a: "1"
---
b: "2"
---
c: "3"

```

`setcfg` outputs to stdout, meaning the results can be piped to a new file or to be included in the results of something like a `kubectl apply` as follows:

**Pipe to file**:
```
$ setcfg -i input.yaml -e dev.yaml > output.yaml
```

**Pipe to kubectl apply**:
```
$ setcfg -i input.yaml -e dev.yaml | kubectl apply -f -
```
