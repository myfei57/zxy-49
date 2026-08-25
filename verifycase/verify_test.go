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

func TestVoScreenPlaylistOrder(t *testing.T) {
	fmt.Println("TestVoScreenPlaylistOrder")
	dir, err := os.MkdirTemp("", "vo-screen")
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
	first, err := screenSvc.Tick("s1", time.Date(2026, 8, 25, 12, 5, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if first.Name != "clip-b" {
		t.Fatalf("first tick must advance by playlist sequence, got %s", first.Name)
	}
	second, err := screenSvc.Tick("s1", time.Date(2026, 8, 25, 12, 6, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if second.Name != "clip-c" {
		t.Fatalf("second tick must advance by playlist sequence, got %s", second.Name)
	}
}
