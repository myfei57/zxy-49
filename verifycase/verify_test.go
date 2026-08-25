package verifycase

import (
	"fmt"
	"os"
	"testing"

	"venueops/internal/camera"
	"venueops/internal/store"
)

func TestVoCameraDayNightThreshold(t *testing.T) {
	fmt.Println("TestVoCameraDayNightThreshold")
	dir, err := os.MkdirTemp("", "vo-daynight")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.SaveSeason(store.Season{Name: "summer", NightHour: 19, NightMin: 30}); err != nil {
		t.Fatal(err)
	}
	cameraSvc := camera.NewService(st, nil)
	if err := cameraSvc.Create("c1", "h1", "Cam One"); err != nil {
		t.Fatal(err)
	}
	if err := st.SaveSeason(store.Season{Name: "winter", NightHour: 17, NightMin: 30}); err != nil {
		t.Fatal(err)
	}
	mode, err := cameraSvc.DayNightSwitch("c1", 18, 0)
	if err != nil {
		t.Fatal(err)
	}
	if mode != "night" {
		t.Fatalf("dusk must switch to night mode using the current threshold, got %s", mode)
	}
}
