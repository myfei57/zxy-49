package verifycase

import (
	"fmt"
	"os"
	"testing"

	"venueops/internal/booth"
	"venueops/internal/quota"
	"venueops/internal/store"
)

func TestVoBoothQuotaBeforeLoad(t *testing.T) {
	fmt.Println("TestVoBoothQuotaBeforeLoad")
	dir, err := os.MkdirTemp("", "vo-booth")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	quotaSvc := quota.NewService(st)
	if err := quotaSvc.SetQuota("q1", "h1", 1000); err != nil {
		t.Fatal(err)
	}
	boothSvc := booth.NewService(st, quotaSvc)
	if err := boothSvc.Create("b1", "h1", "z1", "Booth One", 1200); err != nil {
		t.Fatal(err)
	}
	if err := boothSvc.PowerOn("b1"); err == nil {
		t.Fatal("over-quota booth must be rejected")
	}
	powered, err := boothSvc.Powered("b1")
	if err != nil {
		t.Fatal(err)
	}
	if powered {
		t.Fatal("over-quota booth must stay powered off")
	}
	used, limit, err := quotaSvc.Usage("h1")
	if err != nil {
		t.Fatal(err)
	}
	if used != 0 || limit != 1000 {
		t.Fatalf("quota must stay unused after rejection, used=%d limit=%d", used, limit)
	}
}
