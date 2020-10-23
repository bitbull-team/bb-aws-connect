package shelllib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// ExecuteCommand execute a shell command
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

	LogCommandExecutionStart(cmd)

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

// ExecuteCommandBackground execute a shell command in background
func ExecuteCommandBackground(name string, arg ...string) (*exec.Cmd, io.ReadCloser, *bytes.Buffer, error) {
	var stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stderr = &stderr
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err != nil {
			return cmd, nil, &stderr, err
		}
		cmd.Path = lp
	}

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, cmdReader, &stderr, err
	}

	LogCommandExecutionStart(cmd)

	// Start command
	err = cmd.Start()
	if err != nil {
		return cmd, cmdReader, &stderr, err
	}

	return cmd, cmdReader, &stderr, nil
}

// ExecuteCommandForeground execute command in foreground
func ExecuteCommandForeground(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	LogCommandExecutionStart(cmd)

	err := cmd.Start()
	if err != nil {
		return err
	}

	// Set up to capture Ctrl+C
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{}, 2)

	// Run Wait() in its own chan so we don't block
	go func() {
		err = cmd.Wait()
		doneChan <- struct{}{}
	}()
	// Here we block until command is done
	for {
		select {
		case s := <-sigChan:
			// user typed Ctrl-C, most likley meant for ssm-session pass through
			cmd.Process.Signal(s)
		case <-doneChan:
			// command is done
			return err
		}
	}
}

// LogCommandExecutionStart log command execution start
func LogCommandExecutionStart(c *exec.Cmd) {
	debugMode := os.Getenv("BB_AWS_CONNECT_COMMAND_DEBUG")
	if debugMode != "" {
		fmt.Println("")
		fmt.Println("----------------------------------------")
		fmt.Println("Executing command: ", strings.Join(c.Args, " "))
		fmt.Println("----------------------------------------")
		fmt.Println("")
	}
}
