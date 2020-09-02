# Bitbull CLI

This project is CLI used to collect together repeatable commands and procedures into a single high-level command.

## Install CLI

Download last archive package version from [releases page](https://github.com/bitbull-team/bb-cli/releases):

* Windows: bb-cli_VERSION_Windows_x86_64.zip
* Mac: bb-cli_VERSION_Darwin_x86_64.tar.gz
* Linux: bb-cli_VERSION_Linux_x86_64.tar.gz

Unpack it and copy `bb-cli` into one of your executable paths, for example, for Mac and Linux users:
```bash
tar -czvf bb-cli_*.tar.gz
sudo mv bb-cli /usr/local/bin/bb-cli
rm bb-cli_*.tar.gz
```

### For Linux Users

You can also install CLI from deb or rpm package downloading from releases page:

* bb-cli_VERSION_linux_amd64.deb
* bb-cli_VERSION_linux_amd64.rpm

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

## Debug and troubleshooting

If you want to see executed command (both in background and foreground) set `BB_CLI_DEBUG` environment variable with any value:
```
export BB_CLI_COMMAND_DEBUG=yes
```
then when you execute a command the CLI will print the full command parameters:
```
aws ssm:tunnel --host db.prod.internal --port 3306 --key /my/key/path --username ec2-user --local-port 3306

----------------------------------------
Executing command:  aws ssm start-session --profile myprofile --region eu-west-1 --target i-0cd15458284749f64 --document-name AWS-StartPortForwardingSession --parameters portNumber=22,localPortNumber=59392
----------------------------------------

SSH tunnel to remote instance opened on local port: 59392
Tunnel to remote db.prod.internal:3306 is available on local port: 3306

----------------------------------------
Executing command:  ssh -i /my/key/path -o StrictHostKeyChecking=no -p 59392 ec2-user@localhost -L 3306:db.prod.internal:3306 -T -q
----------------------------------------
```
