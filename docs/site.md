# Site

The `r2pm` *site* is a directory that is only managed by the software.
The user should never touch its contents manually.  
The site is initialized using the `r2pm init` command.


## Location

For each supported operating system, by order of preference:

- Linux:
  - `${XDG_DATA_HOME}/RadareOrg/r2pm` if `$XDG_DATA_HOME` is defined;
  - `${HOME}/.local/share/RadareOrg/r2pm` otherwise
- BSD (including macOS): `${HOME}/Library/RadareOrg/r2pm`
- Windows
  - `${APPDATA}/RadareOrg/r2pm` if `$APPDATA` is defined
  - `${HOMEPATH}/RadareOrg/r2pm` otherwise

## Contents

```
$R2PM_SITE
├── installed/
│   └── pkg1.yaml
│   └── pkg2.yaml
│   └── pkg-from-cli.yaml
├── r2pm-db/
    └── db/
        └── pkg1.yaml
        └── pkg2.yaml
```