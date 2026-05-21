package save

import (
	"path/filepath"
	"testing"

	"kestrelpost/internal/ending"
)

func TestSaveLoadSnapshot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.json")
	in := Snapshot{
		State: ending.RunState{
			Night: 7,
			Fuel:  55,
		},
		Logbook:     []string{"Night 7: test"},
		UnlockCount: 2,
	}
	if err := Save(path, in); err != nil {
		t.Fatalf("save: %v", err)
	}
	out, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if out.State.Night != 7 || out.UnlockCount != 2 {
		t.Fatalf("loaded snapshot mismatch: %+v", out)
	}
}
