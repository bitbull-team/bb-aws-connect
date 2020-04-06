package dockerlib

import (
	"errors"
	"filesystemlib"
	"path"
	"shelllib"
)

// LoginToRepository will login to remote Docker repository
func LoginToRepository(repository string, username string, password string) error {
	args := []string{
		"login",
		"--username=" + username,
		"--password=" + password,
		repository,
	}
	err := shelllib.ExecuteCommand("docker", args...)
	if err != nil {
		return errors.New("Error executing Docker build command: " + err.Error())
	}

	return nil
}

// BuildImage build a Docker image
func BuildImage(contextPath string, tag string, buildArgs []string) error {
	dockerFilePath := path.Join(contextPath, "Dockerfile")
	if filesystemlib.FileExist(dockerFilePath) == false {
		return errors.New("cannot find Dockerfile " + dockerFilePath)
	}

	args := []string{
		"build",
		".",
	}
	for _, buildArg := range buildArgs {
		args = append(args, "--build-arg="+buildArg)
	}
	if tag != "" {
		args = append(args, "--tag="+tag)
	}

	err := shelllib.ExecuteCommand("docker", args...)
	if err != nil {
		return errors.New("Error executing Docker build command: " + err.Error())
	}

	return nil
}
