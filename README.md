# cgrps

**!!! THIS IS STILL A WORK IN PROGRESS !!!**

`cgrps` is a set of commands for checking cgroups.

## Usage

`cgrps` is supposed to be used with [peco](https://github.com/peco/peco) like following command,

```sh
$ cgrps ps $(cgrps ls | peco)
```

## Commands

### `cgrps ls`

list cgroups.

### `cgrps ps [CGROUP]`

report a snapshot of the current cgroups processes.

### `cgrps stat`

show current cgroups stats.
