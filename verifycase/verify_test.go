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

func TestVoZonePartitionFresh(t *testing.T) {
	fmt.Println("TestVoZonePartitionFresh")
	dir, err := os.MkdirTemp("", "vo-zone")
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
	if err := zoneSvc.CreateZone("z-east", "h1", "v1", "East"); err != nil {
		t.Fatal(err)
	}
	if err := zoneSvc.CreateZone("z-west", "h1", "v1", "West"); err != nil {
		t.Fatal(err)
	}
	guardSvc := guard.NewService(st, zoneSvc, hallSvc, auditSvc)
	if err := zoneSvc.Split("h1", "east", "west"); err != nil {
		t.Fatal(err)
	}
	target, err := guardSvc.RouteAlarm("z-west")
	if err != nil {
		t.Fatal(err)
	}
	if target != "west" {
		t.Fatalf("alarm routing must follow the refreshed partition, got %s", target)
	}
}
