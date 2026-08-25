package verifycase

import (
	"fmt"
	"os"
	"testing"

	"venueops/internal/audit"
	"venueops/internal/domain"
	"venueops/internal/hall"
	"venueops/internal/light"
	"venueops/internal/screen"
	"venueops/internal/store"
	"venueops/internal/zone"
)

func TestVoEvacuationOverridePriority(t *testing.T) {
	fmt.Println("TestVoEvacuationOverridePriority")
	dir, err := os.MkdirTemp("", "vo-evac")
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
	screenSvc := screen.NewService(st)
	if err := screenSvc.Create("s1", "Screen One"); err != nil {
		t.Fatal(err)
	}
	lightSvc := light.NewService(st, auditSvc)
	lightSvc.RegisterScreen("s1")
	lightSvc.SetHandover(screenSvc.HandoverToEmergency)
	scheduled := domain.Scene{ZoneID: "z1", Name: "show-scene", Brightness: 70, Color: "warm"}
	if err := lightSvc.UpdateScene(scheduled, 100); err != nil {
		t.Fatal(err)
	}
	if err := lightSvc.Evacuate("h1", 200); err != nil {
		t.Fatal(err)
	}
	effective, err := lightSvc.EffectiveScene("z1")
	if err != nil {
		t.Fatal(err)
	}
	if !effective.IsEmergency {
		t.Fatalf("evacuation must override the scheduled scene, got %s", effective.Name)
	}
	if effective.Brightness != 100 {
		t.Fatalf("emergency scene must run at full brightness, got %d", effective.Brightness)
	}
	handover, err := screenSvc.HandoverState("s1")
	if err != nil {
		t.Fatal(err)
	}
	if !handover {
		t.Fatal("screens must hand over to emergency during evacuation")
	}
}
