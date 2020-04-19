module dockerlib

go 1.14

require (
	filesystemlib v0.0.0
	shelllib v0.0.0
)

replace (
	filesystemlib => ./../filesystem
	shelllib => ../shell
)
