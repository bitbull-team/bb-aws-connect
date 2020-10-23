## AWS commands

This category contain commands used to interact with AWS services.

- [ECS commands](ecs/README.md)
- [SSM commands](ssm/README.md)

### Infrastructure configurations

EC2 instance should have these tags to be able to filter them using `--env` and `--service` parameters:
```
Environment: stage/test/prod
ServiceType: frontend/varnish/ssr/tool
```
no particular value is required, the field is free to be customized as desired.
