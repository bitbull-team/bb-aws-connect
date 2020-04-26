module app

go 1.14

require (
	applib v0.0.0
	configlib v0.0.0
	filesystemlib v0.0.0
	github.com/urfave/cli/v2 v2.2.0
	gopkg.in/yaml.v2 v2.2.8
	shelllib v0.0.0
)

replace (
	applib => ./../../lib/app
	configlib => ./../../lib/config
	filesystemlib => ./../../lib/filesystem
	shelllib => ./../../lib/shell
)
