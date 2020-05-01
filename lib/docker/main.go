package dockerlib

import (
	"errors"
	"filesystemlib"
	"fmt"
	"path"
	"shelllib"
)

// LoginToRegistry login to remote Docker registry
func LoginToRegistry(registry string, username string, password string) error {
	args := []string{
		"login",
		"--username=" + username,
		"--password=" + password,
		registry,
	}
	err := shelllib.ExecuteCommand("docker", args...)
	if err != nil {
		return err
	}

	return nil
}

// BuildImage build a Docker image
func BuildImage(contextPath string, image string, tag string, buildArgs []string) error {
	dockerFilePath := path.Join(contextPath, "Dockerfile")
	if filesystemlib.FileExist(dockerFilePath) == false {
		return errors.New("Cannot find Dockerfile " + dockerFilePath)
	}

	args := []string{
		"build",
		".",
	}
	for _, buildArg := range buildArgs {
		args = append(args, "--build-arg="+buildArg)
	}
	if tag != "" {
		args = append(args, fmt.Sprintf("--tag=%s:%s", image, tag))
	}

	err := shelllib.ExecuteCommand("docker", args...)
	if err != nil {
		return err
	}

	return nil
}

// PushImage push a Docker image
func PushImage(image string, tag string) error {
	args := []string{
		"push",
		fmt.Sprintf("%s:%s", image, tag),
	}
	err := shelllib.ExecuteCommand("docker", args...)
	if err != nil {
		return err
	}

	return nil
}
