package save

import (
	"encoding/json"
	"os"
	"path/filepath"

	"kestrelpost/internal/ending"
)

type Snapshot struct {
	State       ending.RunState `json:"state"`
	Logbook     []string        `json:"logbook"`
	UnlockCount int             `json:"unlock_count"`
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".kestrel-save.json"
	}
	return filepath.Join(home, ".kestrel", "save.json")
}

func Save(path string, s Snapshot) error {
	if path == "" {
		path = DefaultPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o600)
}

func Load(path string) (Snapshot, error) {
	if path == "" {
		path = DefaultPath()
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var out Snapshot
	if err := json.Unmarshal(raw, &out); err != nil {
		return Snapshot{}, err
	}
	return out, nil
}
