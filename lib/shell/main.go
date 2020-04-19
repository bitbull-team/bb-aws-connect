package shelllib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
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

// ExecuteCommandForeground will execute command in foreground
func ExecuteCommandForeground(name string, arg ...string) error {
	rawCmd := exec.Command(name, arg...)
	rawCmd.Stdin = os.Stdin
	rawCmd.Stdout = os.Stdout
	rawCmd.Stderr = os.Stderr

	err := rawCmd.Start()
	if err != nil {
		return err
	}

	// Set up to capture Ctrl+C
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{}, 2)

	// Run Wait() in its own chan so we don't block
	go func() {
		err = rawCmd.Wait()
		doneChan <- struct{}{}
	}()
	// Here we block until command is done
	for {
		select {
		case s := <-sigChan:
			// user typed Ctrl-C, most likley meant for ssm-session pass through
			rawCmd.Process.Signal(s)
		case <-doneChan:
			// command is done
			return err
		}
	}
}
