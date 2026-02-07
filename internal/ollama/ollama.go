package ollama

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func CheckAvailable(ctx context.Context) error {
	// Ollama default: http://localhost:11434
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:11434/", nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot reach ollama at localhost:11434 (%v)", err)
	}
	_ = resp.Body.Close()
	// any HTTP response means it's reachable; we don't enforce status code here
	return nil
}