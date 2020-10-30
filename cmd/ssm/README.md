## SSM commands

This category contain commands used to interact with AWS SSM service.

### Application configuration

Configuration file (by default `.bb-aws-connect.yml`) ha the following configurations:
```yml
profile: default # AWS CLI profile name
region: "eu-west-1" # AWS region
ssm:
   shell: /bin/bash # Shell to open when connecting to instance
   user: root # User to open session with
   cwd: /tmp/ # Command current working directory
```

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
   bb-aws-connect ssm connect - Connect to an EC2 instance using SSM session

USAGE:
   bb-aws-connect ssm connect [command options] [arguments...]

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
```
$ bb-aws-connect ssm connect

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
```
$ bb-aws-connect ssm connect -s cron

? Select an instance: 

  Instace ID          	IP address     	Environment	ServiceType
  [Use arrows to move, type to filter]
> i-0c9916aa684e69638 	172.31.2.57    	stage   	cron
  i-0d6634056d1f36a8a 	172.31.2.93    	test    	cron
  i-074c875d553eca71a 	172.31.2.11    	prod    	cron
  i-0473d16ef7613b580 	172.31.3.186   	redesign	cron
```

also filter by "Environment" tag, is only one instance is found will be auto-selected
```
$ bb-aws-connect ssm connect -e stage -s cron
Instace auto selected: i-0c9916aa684e69638

Starting session with SessionId: botocore-session-1592822987-0fd0966b96e9fde39
root@ip-172-31-2-57:/
```

connect to specific EC2 instance
```
$ bb-aws-connect ssm connect -i i-03edd0f3d32f34b58
```

#### Execute a command to remote EC2 instance

```
NAME:
   bb-aws-connect ssm run - Run command to EC2 instances using a SSM command

USAGE:
   bb-aws-connect ssm run [command options] [command to execute]

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
$ bb-aws-connect ssm run -e test

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
$ bb-aws-connect ssm run -e test -s cron "date"
Instace auto selected:  i-0d6634056d1f36a8a
Waiting for command id  456049a5-868e-4320-9c77-3cbfe362dbcd ..
All commands ends successfully!

--------------------------------
i-0d6634056d1f36a8a 	Success   

Mon Jun 22 11:02:53 UTC 2020

--------------------------------
```

execute "date" command on all instance matching the filter (using `-a` or `--auto-select` parameters)
```
$ bb-aws-connect ssm run -a -e test "date"

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

execute a bash script (for example `my-script.sh`)
```
#!/bin/bash

date
```
on a specific instances (or selected from the list or auto-selected)
```
$ bb-aws-connect ssm run -i i-0f28b75fc851bb344 -i i-0c9916aa684e69638 --file ./my-script.sh

Waiting for command id  745e1051-291a-4445-bc8b-802d5b553b6e ..
All commands ends successfully!

--------------------------------
i-0f28b75fc851bb344 	Success   

Mon Jun 22 12:14:09 UTC 2020

--------------------------------
i-0c9916aa684e69638 	Success   

Mon Jun 22 12:14:09 UTC 2020

--------------------------------
```

#### Execute a SSM Document to remote EC2 instance

```
NAME:
   bb-aws-connect ssm run-document - Run a SSM document to EC2 instances

USAGE:
   bb-aws-connect ssm run-document [command options] [arguments...]

OPTIONS:
   --profile value, -p value    AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
   --region value, -r value     AWS region [$AWS_DEFAULT_REGION]
   --service value, -s value    Service Type (example: bastion, frontend, varnish)
   --env value, -e value        Environment (example: test, stage, prod)
   --instance value, -i value   Instace ID (example: i-xxxxxxxxxxxxxxxxx)
   --auto-select, -a            Automatically select all instance listed without asking (default: false)
   --self, -o                   SSM filter document with owner Self (default: false)
   --private, -t                SSM filter document with owner Private (default: false)
   --document value, -d value   SSM document name
   --parameter value, -m value  SSM document parameter
   --help, -h                   show help (default: false)
```

This command allow you to execute an already deployed [Systems Manager Documents](https://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-ssm-docs.html) to a remote EC2 instance and return the output.

**Usage Example**

use interactive mode to select a SSM Document
```
? Select a document:   [Use arrows to move, type to filter]
> AWS-ASGEnterStandby
  AWS-ASGExitStandby
  AWS-ApplyAnsiblePlaybooks
  AWS-ApplyChefRecipes
  AWS-ApplyDSCMofs
  AWS-ApplyPatchBaseline
  AWS-AttachEBSVolume
  AWS-AttachIAMToInstance
  AWS-ConfigureAWSPackage
  AWS-ConfigureCloudTrailLogging
```

filter only private documents
```
$ bb-aws-connect ssm run-document --private
? Select a document:   [Use arrows to move, type to filter]
> Archive-ZipCompress
  Deploy-Clean
  Archive-TarCompress
  Archive-ZipExtract
  Archive-TarExtract
```

or provide a specific document name
```
$ bb-aws-connect ssm run-document --document test-Dummy
```

document parameters will be asked in interactive mode (it show name and description)
```
? (Required) This is an required parameter. test

? (Optional) This is an optional parameter. 

? (Optional) This has a default value. DefaultValue

? (Optional) This has a default value from Terraform project. TestValue

? (Optional) This parameter has a validation, you can only set it to TestValue. TestValue

? (Optional) This parameter has a validation, you can only write a numeric value. 1
```

then will ask for instance (same as `ssm:run` command), and execute it. At the end of execution the output will be showed
```
? Select an instance: 

  Instace ID            IP address      Environment     ServiceType
 i-0d6634056d1f36a8a    172.31.2.93     test            cron

Waiting for command id  0c265d5f-fd9f-47a3-b434-ba5b152e1472 ..
All commands ends successfully!

--------------------------------
i-0d6634056d1f36a8a     Success   

2020-06-22-12:22:15 value from parameter RequiredParameter is: test
2020-06-22-12:22:15 value from parameter OptionalParameter is: 
2020-06-22-12:22:15 value from parameter ParameterWithDefault is: DefaultValue
2020-06-22-12:22:15 value from parameter ParameterWithDefaultFromTemplate is: TestValue
2020-06-22-12:22:15 value from parameter ParameterWithValidation is: TestValue
2020-06-22-12:22:15 value from parameter ParameterWithValidationRegex is: 1

--------------------------------
```

#### Open an SSH Tunnel

```
NAME:
   bb-aws-connect ssm tunnel - Open a SSM tunnel to a remote host

USAGE:
   bb-aws-connect ssm tunnel [command options] [arguments...]

OPTIONS:
   --profile value, -p value       AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
   --region value, -r value        AWS region [$AWS_DEFAULT_REGION]
   --service value, -s value       Service Type to connect to (example: bastion, frontend, varnish)
   --env value, -e value           Environment (example: test, stage, prod)
   --instance value, -i value      Instace ID (example: i-xxxxxxxxxxxxxxxxx)
   --host value, -o value          Remote host to open tunnel to (example: myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com) (default: "localhost")
   --port value, --rp value        Remote port to open tunnel to (example: 22) (default: "22")
   --local-port value, --lp value  Local port to bind to serve tunnel (example: 2222) (default: "2222")
   --key value, -k value           SSH key to use to connect to instance (example: /path/my-key-pair.pem) (default: "~/.ssh/id_rsa")
   --username value, -u value      SSH username to use to connect to instance (example: ec2-user) (default: "fabio") [$USER]
   --help, -h                      show help (default: false)
```

This command will use [Systems Manager Session](https://aws.amazon.com/it/premiumsupport/knowledge-center/systems-manager-ssh-vpc-resources/) to open an SSH to a specific EC2 instance. 

This command require a valid user created into the EC2 instance (or using the default "ec2-user") and the related SSH key (provided during EC2 creation).

**Usage Example**

open a SSH tunnel specify the user to connect to, (list and filter works as `ssm:connect` command)
```
bb-aws-connect ssm tunnel -s bastion -e test -u fabio.gollinucci
```

an SSH tunnel will be opened locally
```
Instace auto selected: i-07b032553d8f2c0f7
SSH tunnel to remote instance opened on local port: 2222
```

you can now connect to SSH locally
```
ssh myuser@127.0.0.1 -p 2222
```

open a tunnel to a remote RDS instance
```
bb-aws-connect ssm tunnel --host myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com --port 3306 --key instance-key.pem --username ec2-user --local-port 3306
```

instance to connect to will be asked interacvly if not provided and second tunnel (for MySQL port) on the original SSH tunnel:
```
? Select an instance: 

  Instace ID            IP address      Environment     ServiceType
 i-0a9f68b9f8d96d0b1    10.0.1.81       stage   

SSH tunnel to remote instance opened on local port: 56789
Tunnel to remote myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com:3306 is available on local port: 3306
```

you can now connect to MySQL locally (use "127.0.0.1" and not "localhost", otherwise MySQL will try to use the local socket "/var/run/mysqld/mysqld.sock")
```
mysql -h 127.0.0.1 -u myadminuser -p myschema
```
