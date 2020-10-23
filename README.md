# Bitbull AWS Connect CLI

This project is CLI used to collect together repeatable commands and procedures into a single high-level command.

## Install CLI

Download last archive package version from [releases page](https://github.com/bitbull-team/bb-aws-connect/releases):

* Windows: bb-aws-connect_VERSION_Windows_x86_64.zip
* Mac: bb-aws-connect_VERSION_Darwin_x86_64.tar.gz
* Linux: bb-aws-connect_VERSION_Linux_x86_64.tar.gz

Unpack it and copy `bb-aws-connect` into one of your executable paths, for example, for Mac and Linux users:
```bash
tar -czvf bb-aws-connect_*.tar.gz
sudo mv bb-aws-connect /usr/local/bin/bb-aws-connect
rm bb-aws-connect_*.tar.gz
```

### For Linux Users

You can also install CLI from deb or rpm package downloading from releases page:

* bb-aws-connect_1.0.0_linux_amd64.deb
* bb-aws-connect_1.0.0_linux_amd64.rpm

### For Mac Users

Unlock CLI executable file going to "System Preference > Security and Privacy > General" and click on button "open anyway".

## Commands

Check more infos at [commands documentations](cmd/README.md)

## Running commands into different directories

You can execute command into a different directory changing CWD (current working directory) using `--root` flag:
```bash
bb-aws-connect --root ~/Projects/MyProjectA app build
bb-aws-connect --root ~/Projects/MyProjectB app build
```

## Configuration file

This CLI use some options to provide defaults to commands parameters. 
By default CLI will search for `.bb-aws-connect.yml` file into current working directory (if not changed by `--root` flag).

You can override this behaviour using `--config` flag:
```
bb-aws-connect --config /etc/bb-aws-connect.yml app build
```

An example of configuration file can be found into this project root `.bb-aws-connect.yml` file:
```yml
profile: default
region: "eu-west-1"
ecs:
  cluster: default
ssm:
  shell: /bin/bash
  user: root
```

## Wants to contribute?

Here a [contributing guide](CONTRIBUTING.md)

## Debug and troubleshooting

Is possible to enable some verbose logging in order to identify errors and issues.

### Shell commands

If you want to see executed command (both in background and foreground) set `BB_AWS_CONNECT_COMMAND_DEBUG` environment variable with any value:
```
export BB_AWS_CONNECT_COMMAND_DEBUG=yes
```
then when you execute a command the CLI will print the full command parameters:
```
bb-aws-connect ssm tunnel --host db.prod.internal --port 3306 --key /my/key/path --username ec2-user --local-port 3306

----------------------------------------
Executing command:  aws ssm start-session --profile myprofile --region eu-west-1 --target i-0cd15458284749f64 --document-name AWS-StartPortForwardingSession --parameters portNumber=22,localPortNumber=59392
----------------------------------------

SSH tunnel to remote instance opened on local port: 59392
Tunnel to remote db.prod.internal:3306 is available on local port: 3306

----------------------------------------
Executing command:  ssh -i /my/key/path -o StrictHostKeyChecking=no -p 59392 ec2-user@localhost -L 3306:db.prod.internal:3306 -T -q
----------------------------------------
```

### AWS API requests

If you want to see AWS API executed set `BB_AWS_CONNECT_AWS_DEBUG` environment variable with any value:
```
export BB_AWS_CONNECT_AWS_DEBUG=yes
```
then when you execute a command the CLI will print the full command parameters:
```
bb-aws-connect ssm connect -s cron -e test

---[ REQUEST POST-SIGN ]-----------------------------
POST / HTTP/1.1
Host: ec2.eu-west-1.amazonaws.com
User-Agent: aws-sdk-go/1.30.9 (go1.15.2; linux; amd64)
Content-Length: 210
Authorization: AWS4-HMAC-SHA256 Credential=XXXXXXXXXXXXXXXXXX/20201023/eu-west-1/ec2/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-amz-security-token, Signature=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Content-Type: application/x-www-form-urlencoded; charset=utf-8
X-Amz-Date: 20201023T092849Z
X-Amz-Security-Token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Accept-Encoding: gzip

Action=DescribeInstances&Filter.1.Name=instance-state-name&Filter.1.Value.1=running&Filter.2.Name=tag%3AEnvironment&Filter.2.Value.1=test&Filter.3.Name=tag%3AServiceType&Filter.3.Value.1=cron&Version=2016-11-15
-----------------------------------------------------
2020/10/23 11:28:50 DEBUG: Response ec2/DescribeInstances Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK
Transfer-Encoding: chunked
Content-Type: text/xml;charset=UTF-8
Date: Fri, 23 Oct 2020 09:28:49 GMT
Server: AmazonEC2
Vary: accept-encoding
X-Amzn-Requestid: 32631f6f-7c31-442e-8489-b6659fba9601


-----------------------------------------------------
2020/10/23 11:28:50 <?xml version="1.0" encoding="UTF-8"?>
<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
    <requestId>32631f6f-7c31-442e-8489-b6659fba9601</requestId>
    <reservationSet>
        <item>
            <reservationId>r-04bc33g3b153hb3e8</reservationId>
            <ownerId>280215039121</ownerId>
            <groupSet/>
            <instancesSet>
                <item>
                    <instanceId>i-07fe49beh29s7d406</instanceId>
                    <imageId>ami-0d3a49g55e266bee0</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
```
