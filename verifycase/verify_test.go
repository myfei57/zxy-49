package verifycase

import (
	"fmt"
	"os"
	"testing"

	"venueops/internal/audit"
	"venueops/internal/guard"
	"venueops/internal/hall"
	"venueops/internal/store"
	"venueops/internal/zone"
)

func TestVoGuardRouteCurrent(t *testing.T) {
	fmt.Println("TestVoGuardRouteCurrent")
	dir, err := os.MkdirTemp("", "vo-guard")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	auditSvc := audit.NewService(st)
	hallSvc := hall.NewService(st)
	if err := hallSvc.CreateHall("h1", "v1", "Hall"); err != nil {
		t.Fatal(err)
	}
	zoneSvc := zone.NewService(st)
	if err := zoneSvc.CreateZone("z1", "h1", "v1", "Zone One"); err != nil {
		t.Fatal(err)
	}
	if err := hallSvc.BoothChange("h1", "b1", "z1", "Booth One", true); err != nil {
		t.Fatal(err)
	}
	if err := hallSvc.BoothChange("h1", "b2", "z1", "Booth Two", true); err != nil {
		t.Fatal(err)
	}
	guardSvc := guard.NewService(st, zoneSvc, hallSvc, auditSvc)
	if err := hallSvc.BoothChange("h1", "b1", "z1", "Booth One", false); err != nil {
		t.Fatal(err)
	}
	if err := hallSvc.BoothChange("h1", "b3", "z1", "Booth Three", true); err != nil {
		t.Fatal(err)
	}
	if err := guardSvc.RebuildRoute("r1", "z1", "Route One"); err != nil {
		t.Fatal(err)
	}
	checkpoints, err := guardSvc.Checkpoints("r1")
	if err != nil {
		t.Fatal(err)
	}
	foundNew := false
	foundOld := false
	for _, checkpoint := range checkpoints {
		if checkpoint.BoothID == "b3" {
			foundNew = true
		}
		if checkpoint.BoothID == "b1" {
			foundOld = true
		}
	}
	if !foundNew {
		t.Fatal("route must include the newly added booth")
	}
	if foundOld {
		t.Fatal("route must drop the retired booth")
	}
}
