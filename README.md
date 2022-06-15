# ansible-find

A CLI tool to find where ansible variables are defined.

## Install

```
go install github.com/jonasbleyl/ansible-find@latest
```

## Usage

```
A CLI tool to find where ansible variables are defined.

This tool will only search YAML files that reside within the
following directories: [group_vars, host_vars, defaults, vars]

Usage:
  ansible-find VARIABLE [DIRECTORY] [flags]

Flags:
  -l, --files-with-matches   only print the filenames of matching files
  -h, --help                 help for ansible-find
  -r, --regex                use regex to search for variables
  -v, --vault string         ansible vault password file (default ".vault")
```

### Example

```console
$ ansible-find foo_bar 

group_vars/vault.yaml
foo_bar: []

roles/foo/defaults/main.yml
foo_bar:
  - 1
  - 2
```

