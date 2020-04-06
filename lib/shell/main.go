package shelllib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

// ExecuteCommand will execute a shell command
func ExecuteCommand(name string, arg ...string) error {
	var stderr bytes.Buffer
	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
	}
	cmd.Stderr = &stderr
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err != nil {
			return err
		}
		cmd.Path = lp
	}

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
