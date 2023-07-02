<p align="center">
    <img src="internal/static/img/logo.svg" width="128">
    <p align="center">🔖 Self-hosted and lightning-fast bookmark manager</p>
    <p align="center">
      <img alt="GitHub Test Workflow Status" src="https://img.shields.io/github/actions/workflow/status/exler/nurli/tests.yml?branch=main">
      <img alt="MIT License" src="https://img.shields.io/github/license/exler/nurli?color=gold">
    </p>
</p>

## Requirements

* Go >= 1.20

## Usage

### Program usage

```bash
USAGE:
   nurli [global options] command [command options] [arguments...]

COMMANDS:
   version   Show current version
   serve     Run web server
   migrate   Apply changes to the database structure
   bookmark  Manage bookmarks
   tag       Manage tags
   import    Import bookmarks from a file
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --data-dir value   [%NURLI_DATA_DIR%]
   --help, -h        show help
```

## License

Copyright (c) 2023 by ***Kamil Marut***

`Nurli` is under the terms of the [MIT License](https://www.tldrlegal.com/l/mit), following all clarifications stated in the [license file](LICENSE).
