package verifycase

import (
	"fmt"
	"os"
	"testing"
	"time"

	"venueops/internal/domain"
	"venueops/internal/screen"
	"venueops/internal/store"
)

func TestVoScreenNoReplayAfterReset(t *testing.T) {
	fmt.Println("TestVoScreenNoReplayAfterReset")
	dir, err := os.MkdirTemp("", "vo-reset")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	screenSvc := screen.NewService(st)
	if err := screenSvc.Create("s1", "Main"); err != nil {
		t.Fatal(err)
	}
	clips := []domain.Clip{
		{ID: "a", Name: "clip-a", Sponsor: "sp-a", DurationSec: 30},
		{ID: "b", Name: "clip-b", Sponsor: "sp-b", DurationSec: 30},
		{ID: "c", Name: "clip-c", Sponsor: "sp-c", DurationSec: 30},
	}
	if err := screenSvc.SetPlaylist("s1", clips); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 25, 10, 0, 0, 0, time.UTC)
	if _, err := screenSvc.Tick("s1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := screenSvc.Tick("s1", now); err != nil {
		t.Fatal(err)
	}
	if err := screenSvc.Reset("s1", 1000); err != nil {
		t.Fatal(err)
	}
	current, err := screenSvc.Current("s1")
	if err != nil {
		t.Fatal(err)
	}
	if current.ID != "a" {
		t.Fatalf("after reset playback must start from the first clip, got %s", current.ID)
	}
}
