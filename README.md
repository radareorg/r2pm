# r2pm: radare2 package manager

This tool is a cross platform package manager for the reverse engineering
framework radare2.

This tool is still a work in progress.

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
