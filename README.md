# Bitbull CLI

This project is CLI used to collect together repeatable commands and procedures into a single high-level command.

## Install CLI

Download last archive package version from [releases page](https://github.com/bitbull-team/bb-cli/releases):

* Windows: bb-cli_<version>_Windows_x86_64.zip
* Mac: bb-cli_<version>_Darwin_x86_64.tar.gz
* Linux: bb-cli_<version>_Linux_x86_64.tar.gz

Unpack it and copy `bb-cli` into one of your executable paths, for example, for Mac and Linux users:
```bash
tar -czvf bb-cli_*.tar.gz
sudo mv bb-cli /usr/local/bin/bb-cli
rm bb-cli_*.tar.gz
```

### For Linux Users

You can also install CLI from deb or rpm package downloading from releases page:

* bb-cli_<version>_linux_amd64.deb
* bb-cli_<version>_linux_amd64.rpm

### For Mac Users

Unlock CLI executable file going to "System Preference > Security and Privacy > General" and click on button "open anyway".

## Commands

Commands are grouped by first argument:
```bash
bb-cli aws <command> # AWS related commands
bb-cli app <command> # Application manipulation/build related commands
```

Check more infos at [commands documentations](cmd/README.md)

## Running commands into different directories

You can execute command into a different directory changing CWD (current working directory) using `--root` flag:
```bash
bb-cli --root ~/Projects/MyProjectA app build
bb-cli --root ~/Projects/MyProjectB app build
```

## Configuration file

This CLI use some options to provide defaults to commands parameters. 
By default CLI will search for `.bb-cli.yml` file into current working directory (if not changed by `--root` flag).

You can override this behaviour using `--config` flag:
```
bb-cli --config /etc/bb-cli.yml app build
```

An example of configuration file can be found into this project root `.bb-cli.yml` file:
```yml
app:
  type: "go"

aws:
  profile: default
  region: "eu-west-1"
  ecs:
    cluster: default
```

## Wants to contribute?

Here a [contributing guide](CONTRIBUTING.md)
