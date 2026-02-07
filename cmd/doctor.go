package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nishchay7pixels/mockspin/internal/docker"
	"github.com/nishchay7pixels/mockspin/internal/ollama"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment prerequisites (Docker, Ollama)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔎 Checking Docker...")
		if err := docker.EnsureDockerAvailable(cmd.Context()); err != nil {
			return err
		}
		fmt.Println("✅ Docker OK")

		fmt.Println("🔎 Checking Ollama (optional for offline smart data)...")
		if err := ollama.CheckAvailable(cmd.Context()); err != nil {
			fmt.Println("⚠️  Ollama not available:", err)
			fmt.Println("   (You can still run mocks; smart dummy data will be disabled.)")
			return nil
		}
		fmt.Println("✅ Ollama OK")
		return nil
	},
}