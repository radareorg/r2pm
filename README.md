# r2pm: radare2 package manager

This tool is a cross platform package manager for the reverse engineering
framework radare2.

It is a rewrite in Go of the [original Shell r2pm script](https://github.com/radareorg/radare2/blob/master/binr/r2pm/r2pm).

This tool is still a work in progress.

| CI | Badges/URL |
|----------|---------------------------------------------------------------------|
| **GithubCI**  | [![Tests Status](https://github.com/radareorg/r2pm/workflows/Go/badge.svg)](https://github.com/radareorg/r2pm/actions?query=workflow%3AGo)|
| **TravisCI** | [![Build Status](https://travis-ci.org/radareorg/r2pm.svg?branch=master)](https://travis-ci.org/radareorg/r2pm)|
| **Dependabot** |[![Dependabot Enablement](https://api.dependabot.com/badges/status?host=github&repo=radareorg/r2pm)](https://github.com/radareorg/r2pm/pulls?q=is%3Aopen+is%3Apr+label%3Adependencies)|
| **Sourcehut** | [![builds.sr.ht status](https://builds.sr.ht/~xvilka/r2pm.svg)](https://builds.sr.ht/~xvilka/r2pm?)|

## Package example

The official database is available [here](https://github.com/radareorg/r2pm-db).

```yaml
name: r2dec
type: git
repo: https://github.com/wargio/r2dec-js
desc: "[r2-r2pipe-node] an Experimental Decompiler"

install:
  - make -C p

uninstall:
  - make -C p uninstall
```
