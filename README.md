# cgrps [![Build Status](https://travis-ci.org/k1LoW/cgrps.svg?branch=master)](https://travis-ci.org/k1LoW/cgrps)

`cgrps` is a set of commands for checking cgroups.

![cgrps.gif](cgrps.gif)

## Usage

`cgrps` is supposed to be used with [peco](https://github.com/peco/peco) like following command,

```sh
$ cgrps stat $(cgrps ls | peco)
```

or

```sh
$ cgrps ls | grep user.slice | head -1 |  cgrps stat
```

### Use with `ps`

```sh
$ ps u --pid $(cgrps ls | peco | cgrps pids | xargs)
```

### Use with `pidstat`

```sh
$ pidstat -dru -h -p $(cgrps ls | peco | cgrps pids | xargs | tr ' ' ',')
```

### Use with `lsof`

```sh
$ lsof -Pn -i -a -p $(cgrps ls | peco | cgrps pids | xargs | tr ' ' ',')
```

## Commands

### `cgrps ls`

list cgroups.

### `cgrps pids [CGROUP...]`

report a snapshot of the current cgroups pids.

### `cgrps stat [CGROUP]`

show current cgroup stats (`CPU` `MEMORY` `BLKIO` `PIDS`).

## !!!NOTICE!!!

`cgrps` displays cgroups with the same hierarchies together.

If you want to check separately, please use `cgroup-tools (Ubuntu)` or `libcgroup-tools (CentOS)` etc.
