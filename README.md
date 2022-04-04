# Bitbull AWS Connect CLI

This project is CLI to simplify remote connection on AWS resources using System Manager service to remotely open session, tunnel and execute commands.

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

- [ECS commands](cmd/ecs/README.md)
- [SSM commands](cmd/ssm/README.md)

## Project configuration file

This CLI use some options to provide defaults to commands parameters. 
By default CLI will search for `.bb-aws-connect.yml` file into current working directory (if not changed by `--root` flag).

You can override this behaviour using `--config` flag:
```
bb-aws-connect --config /etc/bb-aws-connect.yml ssm connect
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

## Infrastructure configurations

### Tags on resources

EC2 instances should have these tags to be able to filter them using `--env` and `--service` parameters:
```
Environment: prod
ServiceType: frontend
```
no particular value is required, can be customized depending on the use cases, for example:
```
Environment: stage
ServiceType: varnish
```
```
Environment: stage
ServiceType: cron
```
these tags can be applied for a single instance or an AutoScalingGroup. Read more about tagging on [AWS documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html)

### SystemManager agent

EC2 instances must have the SystemManager agent installed and connected, follow the [AWS documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-install-ssm-agent.html) and complete the steps.

### ECS Fargate

ECS Service that use Fargate need to must comply with the use of [ECS Exec](https://docs.aws.amazon.com/AmazonECS/latest/userguide/ecs-exec.html).

### Client IAM permissions

IAM user that execute commands require the following permissions:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ecs:ListClusters",
                "ecs:DescribeClusters",
                "ecs:ListServices",
                "ecs:DescribeServices",
                "ecs:ListTasks",
                "ecs:DescribeTasks",
                "ecs:DescribeContainerInstances",
                "ecs:DescribeTaskDefinition",
                "ecs:DescribeTasks"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ssm:ListCommandInvocations",
                "ssm:ListDocuments",
                "ssm:DescribeDocument"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ssm:SendCommand"
            ],
            "Resource": [
              "arn:aws:ssm:*:*:document/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ssm:SendCommand",
                "ssm:StartSession"
            ],
            "Resource": [
                "arn:aws:ec2:*:*:instance/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ssm:TerminateSession"
            ],
            "Resource": [
                "arn:aws:ssm:*:*:session/${aws:username}-*"
            ]
        }
    ]
}
```

Is possible to restrict session access to instances based on tags:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:SendCommand",
                "ssm:StartSession"
            ],
            "Resource": [
                "arn:aws:ec2:*:*:instance/*"
            ],
            "Condition": {
                "StringLike": {
                    "ssm:resourceTag/ServiceType": [
                        "frontend",
                        "varnish",
                        "cron"
                    ]
                }
            }
        }
    ]
}
```
or only deny access to production instances:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Action": [
                "ssm:StartSession"
            ],
            "Resource": [
                "arn:aws:ec2:*:*:instance/*"
            ],
            "Condition": {
                "StringLike": {
                    "ssm:resourceTag/Environment": [
                        "prod"
                    ]
                }
            }
        }
    ]
}
```

read more about this on [AWS documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/getting-started-restrict-access-examples.html)

## Wants to contribute?

Here the [contributing guide](CONTRIBUTING.md) with some additional tips for debug and local testing.
