package tools

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func ExecCommandWithOutput(command string, arg ...string) ([]byte, []byte, error) {
	cmd := exec.Command(command, arg...)

	readCloserStderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	readCloserStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := io.ReadAll(readCloserStderr)
	if err != nil {
		return nil, nil, err
	}
	stdout, err := io.ReadAll(readCloserStdout)

	err = cmd.Wait()
	if err != nil {
		fmt.Println("stdout: ", string(stdout))
		fmt.Println("stderr: ", string(stderr))
		return nil, nil, errors.New("cmd.Wait(): " + err.Error())
	}

	return stdout, stderr, nil
}

func SaveInfile(path string, bytes []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
