## Contributing to CLI

### Install GO

Install GO and development tools, fore info at https://golang.org/doc/install

### Execute locally

To execute the CLI locally from this repository run:
```bash
go run main.go --help
```

### Build

To build CLI into an executable file run:
```bash
go build -o bb-aws-connect main.go
```

You can be able to execute file directly:
```bash
./bb-aws-connect --help
```

### Install a dev version

If you want to override your system installed CLI with development version. 
Build it (if not already did):
```bash
go build -o bb-aws-connect main.go
```
and then install it into your system:
```bash
sudo cp bb-aws-connect /usr/local/bin/bb-aws-connect
```

## Resources

* [A tour of Go](https://tour.golang.org/list)
* [CLI library](https://github.com/urfave/cli/blob/master/docs/v2/manual.md)
* [Input library](https://github.com/AlecAivazis/survey)
