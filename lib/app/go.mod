module applib

go 1.14

require (
	configlib v0.0.0
	filesystemlib v0.0.0
	shelllib v0.0.0
)

replace (
	configlib => ../config
	filesystemlib => ../filesystem
	shelllib => ../shell
)
