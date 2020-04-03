module bb-cli

go 1.14

require (
	github.com/urfave/cli/v2 v2.2.0
	docker v0.0.0
)

replace (
	docker v0.0.0 => "./cmd/docker"
)
