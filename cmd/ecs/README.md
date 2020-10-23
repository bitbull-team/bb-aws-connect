## ECS commands

This category contain commands used to interact with AWS ECS service.

### Application configuration

Configuration file (by default `.bb-aws-connect.yml`) ha the following configurations:
```yml
profile: default # AWS CLI profile name
region: "eu-west-1" # AWS region
ecs:
  cluster: default # ECS cluster name
```

### Commands

All the commands in this section have the following options:
```
--profile value, -p value   AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
--region value, -r value    AWS region [$AWS_DEFAULT_REGION]
```

These options can be configured using `AWS_PROFILE`, `AWS_DEFAULT_PROFILE`, `AWS_DEFAULT_REGION` environment variables (the same as AWS cli).

#### Open a shell to ECS Task's container

```
NAME:
   bb-aws-connect ecs connect - Connect to an ECS Task container

USAGE:
   bb-aws-connect ecs connect [command options] [arguments...]

OPTIONS:
   --profile value, -p value  AWS profile name [$AWS_PROFILE, $AWS_DEFAULT_PROFILE]
   --region value, -r value   AWS region [$AWS_DEFAULT_REGION]
   --cluster value, -c value  Cluster Name
   --service value, -s value  Service name (example: my-service)
   --task value, -t value     Task ID (example: xxxxxxxxxxxxxxxxxxxx)
   --workdir value, -w value  Docker exec 'workdir' parameters (example: /app)
   --user value, -u value     Docker exec 'user' parameters (example: www-data)
   --command value            Use a custom command as entrypoint (default: "/bin/bash")
   --help, -h                 show help (default: false)
```

This command allow you to open an interactive shell using Systems Manager Session Manager directly into a [ECS](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html) Task container.

**Usage Example**

if cluster name is not provided will be asked interactively
```
$ bb-aws-connect ecs connect

? Select a cluster:  [Use arrows to move, type to filter]
> myproject-prod
  myproject-stage
```

if service name is not provided will be asked interactively
```
? Select a service: 

  Name                                          Desired Running
  [Use arrows to move, type to filter]
> myservice-1                                   1       1     
  myservice-3                                   1       1     
```

if task name is not provided will be asked interactively, if service as only one task will be selected automatically
```
Task auto selected:  arn:aws:ecs:eu-west-1:0000000000:task/myproject-stage/d3833a6674344776897fda5d7012aa55
```

finally will be asked which container you want to connect to
```
? Select a container: 

  Container Name                        Container ID      Status        Image
  [Use arrows to move, type to filter]
> task-1-logger                         8bd8e24d5e76     RUNNING        logger:latest
  task-ic-1-cron                        ca9e7f769b75     RUNNING        000000000.dkr.ecr.eu-west-1.amazonaws.com/myproject/cron:1.0.0
```

shell will be opened with provided configuration
```
$ bb-aws-connect ecs connect --workdir /app --user www-data

...

Starting session with SessionId: botocore-session-1592829669-05415380dd16a4f30
www-data@ca9e7f769b75:/app$ 
```
