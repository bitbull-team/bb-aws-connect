module bb-cli

go 1.14

require (
	github.com/urfave/cli/v2 v2.2.0
	docker v0.0.0
	dockerlib v0.0.0
	filesystemlib v0.0.0
	shelllib v0.0.0
	aws v0.0.0
	awslib v0.0.0
	app v0.0.0
	applib v0.0.0
	configlib v0.0.0
)

replace (
	docker => ./cmd/docker
	aws => ./cmd/aws
	app => ./cmd/app

	awslib => ./pkg/aws
	dockerlib => ./pkg/docker
	filesystemlib => ./pkg/filesystem
	shelllib => ./pkg/shell
	applib => ./pkg/app
	configlib => ./pkg/config
)
