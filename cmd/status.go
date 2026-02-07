package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourname/mockspin/internal/docker"
	"github.com/yourname/mockspin/internal/session"
)

var statusProject string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current mock server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		project := statusProject
		if project == "" {
			project = "default"
		}
		s, err := session.Load(project)
		if err != nil {
			return err
		}
		if s == nil || s.ContainerName == "" {
			fmt.Println("ℹ️  No active session.")
			return nil
		}

		running, err := docker.IsRunning(cmd.Context(), s.ContainerName)
		if err != nil {
			return err
		}

		fmt.Println("Project  :", s.Project)
		fmt.Println("Spec     :", s.SpecPath)
		fmt.Println("URL      :", fmt.Sprintf("http://localhost:%d", s.Port))
		fmt.Println("Container:", s.ContainerName)
		fmt.Println("Running  :", running)
		return nil
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusProject, "project", "", "Project name (default: default)")
}
