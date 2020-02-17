# Packages

Packages are additions to the `r2` software.  
A package can consist in a set of source files that have to be built against the current `r2` installation, or in pre-built binaries.

## Plugin file

A plugin is described by a plugin file.
A plugin file is a text file in the YAML format that contains various metadata about the plugin as well as instructions on how to install the package for each platform.

### Schema

```yaml
---
name: my-package  # must be unique and equals to the file name
version: 1.2.3
description: Some description   

install:

  windows:
    source:
      type: zip
      url: http://a-random.url/zip-archive.zip
    out:
      - path: bin/exe1  # relative to the extracted directory
        type: exe
      - path: lib/mylib
        type: shared-lib

  linux:
    source:
      type: git
      repo: git@github.com:username/project.git
      ref: master  # or a tag
    commands:
      - './configure --prefix {{ .DestPath }}'
      - make
    out:
      - path: bin/exe1
        type: exe
      - path: lib/mylib
        type: shared-lib

  macos:
    source:
      type: git
      repo: git@github.com:username/project.git
      ref: master  # or a tag
    commands:
      - gcc -o exe1 exe1.c -I{{ .R2HeadersPath }} -L{{ .R2LibsPath }}
    out:
      - path: exe1
        type: exe
```

### Template variables

To accomodate those plugins that need to compile against `r2` libraries, `r2pm` will replace the following variables with their value in the commands:

|Name|Description|
|----|-----------|
|`DestPath`|The destination directory for plugins|
|`R2HeadersPath`|The directory where `r2` headers are located|
|`R2LibsPath`|The directory where `r2` libraries are located|
