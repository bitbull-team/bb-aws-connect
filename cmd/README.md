## Commands

This directory contain first level commands, for example:
```bash
bb-cli aws <command>
```
you will find a directory named "aws" with sub commands declarations.

This is used to group commands by categories
```bash
bb-cli aws <command> # AWS related commands
bb-cli app <command> # Application manipulation/build related commands
```

## Categories

* [Application](app/README.md)
* [AWS](aws/README.md)
* [Docker](docker/README.md)
