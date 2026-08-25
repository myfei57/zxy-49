package verifycase

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"venueops/internal/audit"
	"venueops/internal/gate"
	"venueops/internal/store"
)

func TestVoGateModeDurableFirst(t *testing.T) {
	fmt.Println("TestVoGateModeDurableFirst")
	dir, err := os.MkdirTemp("", "vo-gate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	auditSvc := audit.NewService(st)
	gateSvc := gate.NewService(st, auditSvc)
	if err := gateSvc.Create("g1", "z1", "Gate One"); err != nil {
		t.Fatal(err)
	}
	journal := filepath.Join(dir, "journal")
	if err := os.RemoveAll(journal); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(journal, []byte("blocked"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gateSvc.SetMode("g1", "pass", 100); err == nil {
		t.Fatal("SetMode must fail when the mode journal cannot be written")
	}
	mode, err := gateSvc.Mode("g1")
	if err != nil {
		t.Fatal(err)
	}
	if mode != "armed" {
		t.Fatalf("gate must stay armed when the mode is not durable, got %s", mode)
	}
	released, err := gateSvc.Released("g1")
	if err != nil {
		t.Fatal(err)
	}
	if released {
		t.Fatal("gate must not release when the mode is not durable")
	}
	if err := os.RemoveAll(journal); err != nil {
		t.Fatal(err)
	}
	if err := gateSvc.SetMode("g1", "pass", 101); err != nil {
		t.Fatal(err)
	}
	durable, err := gateSvc.DurableMode("g1")
	if err != nil {
		t.Fatal(err)
	}
	if durable != "pass" {
		t.Fatalf("durable mode must be pass, got %s", durable)
	}
}
