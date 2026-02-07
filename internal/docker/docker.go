package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yourname/mockspin/internal/execx"
)

type Mount struct {
	HostPath      string
	ContainerPath string
	ReadOnly      bool
}

type RunSpec struct {
	Image         string
	ContainerName string
	HostPort      int
	ContainerPort int
	Mounts        []Mount
	Args          []string
	AutoRemove    bool
}

func EnsureDockerAvailable(ctx context.Context) error {
	_, err := execx.Run(ctx, "docker", "version")
	return err
}

func Run(ctx context.Context, spec RunSpec, detach bool) (string, error) {
	args := []string{"run"}
	if detach {
		args = append(args, "-d")
	}
	if spec.AutoRemove && !detach {
		args = append(args, "--rm")
	}
	args = append(args, "--name", spec.ContainerName)
	args = append(args, "-p", fmt.Sprintf("%d:%d", spec.HostPort, spec.ContainerPort))

	for _, m := range spec.Mounts {
		v := fmt.Sprintf("%s:%s", m.HostPath, m.ContainerPath)
		if m.ReadOnly {
			v += ":ro"
		}
		args = append(args, "-v", v)
	}

	args = append(args, spec.Image)
	args = append(args, spec.Args...)

	res, err := execx.Run(ctx, "docker", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}

func StopAndRemove(ctx context.Context, containerName string) error {
	// stop (ignore errors), then rm (ignore errors)
	_, _ = execx.Run(ctx, "docker", "stop", containerName)
	_, _ = execx.Run(ctx, "docker", "rm", "-f", containerName)
	return nil
}

func IsRunning(ctx context.Context, containerName string) (bool, error) {
	res, err := execx.Run(ctx, "docker", "inspect", "-f", "{{.State.Running}}", containerName)
	if err != nil {
		// if container doesn't exist => not running
		if strings.Contains(err.Error(), "No such object") {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(res.Stdout) == "true", nil
}

func Logs(ctx context.Context, containerName string) error {
	// Use exec.CommandContext to stream directly to stdout/stderr
	cmd := exec.CommandContext(ctx, "docker", "logs", "-f", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func LogsOnce(ctx context.Context, containerName string) error {
	cmd := exec.CommandContext(ctx, "docker", "logs", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}