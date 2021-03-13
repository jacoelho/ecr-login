package ecr

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func DockerLogin(ctx context.Context, c Credential) error {
	cmd := exec.CommandContext(ctx, "docker",
		"login",
		fmt.Sprintf("--username=%s", c.Username),
		"--password-stdin",
		c.RegistryURL)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, err := io.Copy(stdin, strings.NewReader(c.Password)); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	return cmd.Run()
}
