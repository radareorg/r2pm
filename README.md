# r2pm: radare2 package manager

This tool is a cross platform package manager for the reverse engineering
framework radare2.

It is a rewrite in Go of the [original Shell r2pm script](https://github.com/radare/radare2/blob/master/binr/r2pm/r2pm).

This tool is still a work in progress.

| CI | Badges/URL |
|----------|---------------------------------------------------------------------|
| **GolangCI** 	| https://golangci.com/r/github.com/radareorg/r2pm|
| **TravisCI** 	| [![Build Status](https://travis-ci.org/radareorg/r2pm.svg?branch=master)](https://travis-ci.org/radareorg/r2pm)|
| **Appveyor** 	|[![Build status](https://ci.appveyor.com/api/projects/status/3otiyo19a6hnyog3?svg=true)](https://ci.appveyor.com/project/radare/r2pm)|

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
