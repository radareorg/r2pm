# r2pm: radare2 package manager

This tool is a cross platform package manager for the reverse engineering
framework radare2.

This tool is still a work in progress.

## Package example
```
{
"name": "swf2",
"type": "git",
"repo": "https://github.com/radare/radare2-extras",
"desc": "[r2-bin] SWF/Flash disassembler",
"install": [
  "./configure --prefix=\"${R2PM_PREFIX}\" || exit 1",
  "cd libr/asm/p",
  "${MAKE} clean",
  "${MAKE} asm_swf.${LIBEXT} || exit 1",
  "mkdir -p \"${R2PM_PLUGDIR}\" || exit 1",
  "cp -f asm_swf.${LIBEXT} \"${R2PM_PLUGDIR}\" || exit 1",
  "cd ../../bin/p || exit 1",
  "${MAKE} bin_swf.${LIBEXT} || exit 1",
  "echo cp -f bin_swf.${LIBEXT} \"${R2PM_PLUGDIR}\" || exit 1",
  "cp -f bin_swf.${LIBEXT} \"${R2PM_PLUGDIR}\" || exit 1"
],
"uninstall": [
  "rm -f \"${R2PM_PLUGDIR}\"/*swf*"
]
}
```

