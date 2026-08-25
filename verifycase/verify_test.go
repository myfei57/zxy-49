package verifycase

import (
	"fmt"
	"os"
	"testing"

	"venueops/internal/camera"
	"venueops/internal/domain"
	"venueops/internal/hall"
	"venueops/internal/store"
)

func TestVoCameraPresetRefresh(t *testing.T) {
	fmt.Println("TestVoCameraPresetRefresh")
	dir, err := os.MkdirTemp("", "vo-camera")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	hallSvc := hall.NewService(st)
	if err := hallSvc.CreateHall("h1", "v1", "Hall"); err != nil {
		t.Fatal(err)
	}
	if err := hallSvc.SetBoothPosition("h1", "b1", domain.Point{X: 10, Y: 10}); err != nil {
		t.Fatal(err)
	}
	cameraSvc := camera.NewService(st, hallSvc)
	if err := cameraSvc.Create("c1", "h1", "Cam One"); err != nil {
		t.Fatal(err)
	}
	if err := cameraSvc.SetPreset("c1", "b1", 15, 25); err != nil {
		t.Fatal(err)
	}
	if err := hallSvc.LayoutChange("h1", map[string]domain.Point{"b1": {X: 50, Y: 60}}); err != nil {
		t.Fatal(err)
	}
	if err := cameraSvc.RefreshPresets("c1", 0); err != nil {
		t.Fatal(err)
	}
	preset, err := cameraSvc.PresetFor("c1", "b1")
	if err != nil {
		t.Fatal(err)
	}
	if preset.X != 50 || preset.Y != 60 {
		t.Fatalf("preset must follow the current layout, got (%d,%d)", preset.X, preset.Y)
	}
}
