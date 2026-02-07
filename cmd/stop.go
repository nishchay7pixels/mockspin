package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourname/mockspin/internal/docker"
	"github.com/yourname/mockspin/internal/session"
)

var stopProject string

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running mock server for a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		project := stopProject
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
		if err := docker.StopAndRemove(cmd.Context(), s.ContainerName); err != nil {
			return err
		}
		_ = session.Delete(project)
		fmt.Println("✅ Stopped:", s.ContainerName)
		return nil
	},
}

func init() {
	stopCmd.Flags().StringVar(&stopProject, "project", "", "Project name (default: default)")
}