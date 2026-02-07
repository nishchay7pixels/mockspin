package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/nishchay7pixels/mockspin/internal/util"
)

type Session struct {
	Project       string    `json:"project"`
	SpecPath      string    `json:"spec_path"`
	SpecHash      string    `json:"spec_hash"`
	ContainerName string    `json:"container_name"`
	ContainerID   string    `json:"container_id"`
	Port          int       `json:"port"`
	StartedAt     time.Time `json:"started_at"`
}

func sessionPath(project string) (string, error) {
	dir, err := util.SessionsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, project+".json"), nil
}

func Save(s *Session) error {
	p, err := sessionPath(s.Project)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o600)
}

func Load(project string) (*Session, error) {
	p, err := sessionPath(project)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func Delete(project string) error {
	p, err := sessionPath(project)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}