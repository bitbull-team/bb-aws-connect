module bb-cli

go 1.14

require (
	github.com/urfave/cli/v2 v2.2.0
	docker v0.0.0
	dockerlib v0.0.0
	filesystemlib v0.0.0
	shelllib v0.0.0
)

replace (
	docker => ./cmd/docker
	dockerlib => ./lib/docker
	filesystemlib => ./lib/filesystem
	shelllib => ./lib/shell
)
