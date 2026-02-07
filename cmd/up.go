package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourname/mockspin/internal/docker"
	"github.com/yourname/mockspin/internal/ollama"
	"github.com/yourname/mockspin/internal/session"
	"github.com/yourname/mockspin/internal/util"
)

var (
	upSpec    string
	upPort    int
	upProject string
	upDetach  bool
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start a local mock server container from an OpenAPI spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		if upSpec == "" {
			return fmt.Errorf("--spec is required")
		}

		absSpec, err := filepath.Abs(upSpec)
		if err != nil {
			return err
		}
		if _, err := os.Stat(absSpec); err != nil {
			return fmt.Errorf("spec not found: %s", absSpec)
		}

		// Preflight checks
		if err := docker.EnsureDockerAvailable(cmd.Context()); err != nil {
			return err
		}
		// Ollama is optional (offline AI). We'll check and report but not fail.
		_ = ollama.CheckAvailable(cmd.Context())

		// Create a stable container name per project+spec hash
		specHash, err := hashFile(absSpec)
		if err != nil {
			return err
		}
		project := upProject
		if project == "" {
			project = "default"
		}
		containerName := fmt.Sprintf("mockspin_%s_%s", project, specHash[:8])

		// If a session exists for this project, stop it first (clean behavior for MVP)
		if s, _ := session.Load(project); s != nil && s.ContainerName != "" {
			_ = docker.StopAndRemove(cmd.Context(), s.ContainerName)
			_ = session.Delete(project)
		}

		// Start Prism container (MVP)
		// Prism: stoplight/prism, mock mode, mount spec
		// Example: docker run --rm --name <name> -p 4010:4010 -v <spec>:/tmp/openapi.yaml stoplight/prism:4 mock -h 0.0.0.0 /tmp/openapi.yaml
		run := docker.RunSpec{
			Image:         "stoplight/prism:4",
			ContainerName: containerName,
			HostPort:      upPort,
			ContainerPort: 4010,     // ✅ Prism listens on 4010 inside container
			AutoRemove:    false,    // ✅ keep container so logs can be read
			Mounts: []docker.Mount{
				{HostPath: absSpec, ContainerPath: "/tmp/openapi.yaml", ReadOnly: true},
			},
			Args: []string{"mock", "-h", "0.0.0.0", "/tmp/openapi.yaml"},
		}

		containerID, err := docker.Run(cmd.Context(), run, true) // ✅ always -d

		running, _ := docker.IsRunning(cmd.Context(), containerName)
		if !running {
		_ = docker.LogsOnce(cmd.Context(), containerName) // prints last logs
		_ = docker.StopAndRemove(cmd.Context(), containerName)
		_ = session.Delete(project)
		return fmt.Errorf("prism exited immediately (see logs above)")
		}

		if err != nil {
			return err
		}

		// Save session
		s := &session.Session{
			Project:       project,
			SpecPath:      absSpec,
			SpecHash:      specHash,
			ContainerName: containerName,
			ContainerID:   containerID,
			Port:          upPort,
			StartedAt:     time.Now(),
		}
		if err := session.Save(s); err != nil {
			// best effort cleanup if saving fails
			_ = docker.StopAndRemove(cmd.Context(), containerName)
			return err
		}

		fmt.Printf("✅ Mock server running\n")
		fmt.Printf("   Project : %s\n", project)
		fmt.Printf("   Spec    : %s\n", absSpec)
		fmt.Printf("   URL     : http://localhost:%d\n", upPort)
		fmt.Printf("   Container: %s\n", containerName)

		if upDetach {
			fmt.Println("ℹ️  Detached mode. Use `mockspin status` or `mockspin stop`.")
			return nil
		}

		// Foreground mode: attach logs, and stop on exit signals
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Println("\n🛑 Stopping...")
			cancel()
		}()

		// Stream logs until ctx cancelled
		_ = docker.Logs(ctx, containerName)

		// Cleanup on exit
		_ = docker.StopAndRemove(context.Background(), containerName)
		_ = session.Delete(project)
		fmt.Println("✅ Stopped.")
		return nil
	},
}

func init() {
	upCmd.Flags().StringVarP(&upSpec, "spec", "s", "", "Path to OpenAPI spec (yaml/json)")
	upCmd.Flags().IntVarP(&upPort, "port", "p", 4010, "Local port to expose")
	upCmd.Flags().StringVar(&upProject, "project", "", "Project name (session namespace)")
	upCmd.Flags().BoolVar(&upDetach, "detach", false, "Run container in background")
}

func hashFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// used by other packages in later phases
func ensureAppDirs() error {
	_, err := util.SessionsDir()
	return err
}