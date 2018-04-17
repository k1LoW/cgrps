# cgrps

`cgrps` is a set of commands for checking cgroups.

## Usage

`cgrps` is supposed to be used with [peco](https://github.com/peco/peco) like following command,

```sh
$ cgrps ps $(cgrps ls | peco)
```

or

```sh
$ cgrps ls | grep user.slice | head -1 |  cgrps ps
```

## Commands

### `cgrps ls`

list cgroups.

### `cgrps ps [CGROUP]`

report a snapshot of the current cgroups processes.

### `cgrps stat [CGROUP]`

show current cgroups stats (`CPU` `MEMORY` `BLKIO`).
