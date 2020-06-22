## AWS commands

This category contain commands used to interact with AWS services.

### Application configuration

Configuration file (by default `.bb-cli.yml`) ha the following configurations:
```yml
aws:
  profile: default # AWS CLI profile name
  region: "eu-west-1" # AWS region

  ecs:
    cluster: default # ECS cluster name
```

### Infrastructure configurations

EC2 instance should have these tags to be able to filter them using `--env` and `--service` parameters:
```
Environment: stage/test/prod
ServiceType: frontend/varnish/ssr/tool
```
no particular value is required, the field is free to be customized as desired.

### Commands

All the commands in this section have the following options:
```
--profile value, -p value   AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
--region value, -r value    AWS region [$AWS_DEFAULT_REGION]
```

These options can be configured using `AWS_PROFILE`, `AWS_DEFAULT_PROFILE`, `AWS_DEFAULT_REGION` environment variables (the same as AWS cli).

#### Connect to remote EC2 instance

```
NAME:
   bb-cli aws ssm:connect - Connect to an EC2 instance using SSM session

USAGE:
   bb-cli aws ssm:connect [command options] [arguments...]

OPTIONS:
   --profile value, -p value   AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
   --region value, -r value    AWS region [$AWS_DEFAULT_REGION]
   --service value, -s value   Service Type (example: bastion, frontend, varnish)
   --env value, -e value       Environment (example: test, stage, prod)
   --instance value, -i value  Instace ID (example: i-xxxxxxxxxxxxxxxxx)
   --cwd value                 Current working directory (example: /var/www/) (default: "/")
   --user value                User to use in the session (default: "root")
   --shell value               Shell used in session (default: "/bin/bash")
   --command value             Use a custom command as entrypoint
   --help, -h                  show help (default: false)
```

This command allow you to connect to a remote EC2 instance using [Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html).

**Usage Example**

use interactive mode, list can be filter just typing something (search will be performed on all columns)
```bash
$ bb-cli aws ssm:connect

? Select an instance: 

  Instace ID          	IP address     	Environment	ServiceType
  [Use arrows to move, type to filter]
> i-0f8061bcb996a0a1b 	172.31.2.206   	prod    	varnish
  i-0c9916aa684e69638 	172.31.2.57    	stage   	cron
  i-0d6634056d1f36a8a 	172.31.2.93    	test    	cron
  i-074c875d553eca71a 	172.31.2.11    	prod    	cron
  i-046960e357e6e52e2 	172.31.2.178   	prod    	frontend
  i-09aa1106e7b67fb42 	172.31.2.68    	prod    	frontend
  i-055675b814d0d51b8 	172.31.2.245   	prod    	frontend
  i-048e13ab5e9f135f0 	172.31.2.22    	prod    	frontend
  i-091de6392bfa4420b 	172.31.2.190   	prod    	frontend
  i-023c435e2aed08b46 	172.31.2.174   	prod    	frontend
  i-0dfa64cb5f4f44277 	172.31.2.200   	prod    	frontend
  i-06ac8583013436f99 	172.31.3.119   	redesign	ssr
  i-04fedf58b8155c708 	172.31.3.78    	stage   	varnish
```

pre-filter service type checking "ServiceType" tag value
```bash
$ bb-cli aws ssm:connect -s cron

? Select an instance: 

  Instace ID          	IP address     	Environment	ServiceType
  [Use arrows to move, type to filter]
> i-0c9916aa684e69638 	172.31.2.57    	stage   	cron
  i-0d6634056d1f36a8a 	172.31.2.93    	test    	cron
  i-074c875d553eca71a 	172.31.2.11    	prod    	cron
  i-0473d16ef7613b580 	172.31.3.186   	redesign	cron
```

also filter by "Environment" tag, is only one instance is found will be auto-selected
```bash
$ bb-cli aws ssm:connect -e stage -s cron
Instace auto selected: i-0c9916aa684e69638

Starting session with SessionId: botocore-session-1592822987-0fd0966b96e9fde39
root@ip-172-31-2-57:/# 
```

connect to specific EC2 instance
```bash
$ bb-cli aws ssm:connect -i i-03edd0f3d32f34b58
```

#### Execute a command to remote EC2 instance

```
NAME:
   bb-cli aws ssm:run - Run command to EC2 instances using a SSM command

USAGE:
   bb-cli aws ssm:run [command options] [command to execute]

OPTIONS:
   --profile value, -p value   AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
   --region value, -r value    AWS region [$AWS_DEFAULT_REGION]
   --service value, -s value   Service Type (example: bastion, frontend, varnish)
   --env value, -e value       Environment (example: test, stage, prod)
   --instance value, -i value  Instace ID (example: i-xxxxxxxxxxxxxxxxx)
   --file value                Script file path to execute (example: ./my-script.sh)
   --auto-select, -a           Automatically select all instance listed without asking (default: false)
   --help, -h                  show help (default: false)
```

This command allow you to execute a single command to a remote EC2 instance using [Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html) and return the output.

**Usage Example**

use interactive mode, multiple instance can be selected (command will be executed on all instance asynchronously)
```
$ bb-cli aws ssm:run -e test

? Select an instance: 

  Instace ID            IP address      Environment     ServiceType
  [Use arrows to move, space to select, type to filter]
> [ ]  i-0d6634056d1f36a8a      172.31.2.93     test            cron
  [x]  i-0b8be6c25d48949f9      172.31.3.56     test            ssr
  [x]  i-07b032553d8f2c0f7      172.31.4.200    test            bastion
  [ ]  i-0f28b75fc851bb344      172.31.1.30     test            frontend
  [ ]  i-0051e23402345c68a      172.31.1.96     test            ssr
```

is not command provided by arguments it will be asked interactively (multiple command can be provided with multiple line, as a bash script)
```
? Select an instance: 

  Instace ID            IP address      Environment     ServiceType
 i-0b8be6c25d48949f9    172.31.3.56     test            ssr, i-07b032553d8f2c0f7        172.31.4.200    test            bastion

? Type command to execute:  [Enter 2 empty lines to finish]
date
```

when command ends the output will be showed for each instances
```
Waiting for command id  8ef3a896-50e6-4075-a83b-a2671b614a9a ..
All commands ends successfully!

--------------------------------
i-0b8be6c25d48949f9 	Success   

Mon Jun 22 11:02:00 UTC 2020

--------------------------------
i-07b032553d8f2c0f7 	Success   

Mon Jun 22 11:02:00 UTC 2020

--------------------------------
```

one liner command for a single instance
```
$ bb-cli aws ssm:run -e test -s cron "date"
Instace auto selected:  i-0d6634056d1f36a8a
Waiting for command id  456049a5-868e-4320-9c77-3cbfe362dbcd ..
All commands ends successfully!

--------------------------------
i-0d6634056d1f36a8a 	Success   

Mon Jun 22 11:02:53 UTC 2020

--------------------------------
```

execute on all instance matching the filter (using `-a` or `--auto-select` parameters)
```
$ bb-cli aws ssm:run -a -e test "date"

Waiting for command id  d5ef418f-5e9d-432d-af73-a8941df6fe4f ..
All commands ends successfully!

--------------------------------
i-0f28b75fc851bb344 	Success   

Mon Jun 22 11:03:30 UTC 2020

--------------------------------
i-0d6634056d1f36a8a 	Success   

Mon Jun 22 11:03:30 UTC 2020

--------------------------------
```
