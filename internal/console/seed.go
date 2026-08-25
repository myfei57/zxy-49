package console

import (
	"time"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Server) Seed() error {
	if !s.nsSvc.HasVenue("v-main") {
		if err := s.nsSvc.CreateVenue("v-main", "国家会展中心"); err != nil {
			return err
		}
	}
	for _, hall := range []struct {
		id   string
		name string
	}{{"h-a", "A 馆"}, {"h-b", "B 馆"}} {
		if !s.hallSvc.HasHall(hall.id) {
			if err := s.hallSvc.CreateHall(hall.id, "v-main", hall.name); err != nil {
				return err
			}
		}
	}
	zones := []struct {
		id   string
		hall string
		name string
	}{
		{"z-a1", "h-a", "A1 展区"},
		{"z-a2", "h-a", "A2 展区"},
		{"z-b1", "h-b", "B1 展区"},
		{"z-b2", "h-b", "B2 展区"},
	}
	for _, zone := range zones {
		if !s.zoneSvc.HasZone(zone.id) {
			if err := s.zoneSvc.CreateZone(zone.id, zone.hall, "v-main", zone.name); err != nil {
				return err
			}
		}
		if err := s.hallSvc.AddZone(zone.hall, zone.id); err != nil {
			return err
		}
		if err := s.nsSvc.AttachZone("v-main", zone.id); err != nil {
			return err
		}
	}
	booths := []struct {
		id     string
		hall   string
		zone   string
		name   string
		watts  int
		posX   int
		posY   int
	}{
		{"b-101", "h-a", "z-a1", "101 展位", 800, 10, 10},
		{"b-102", "h-a", "z-a1", "102 展位", 600, 30, 10},
		{"b-201", "h-b", "z-b1", "201 展位", 1200, 10, 20},
		{"b-202", "h-b", "z-b2", "202 展位", 900, 30, 20},
	}
	for _, booth := range booths {
		if !s.boothSvc.HasBooth(booth.id) {
			if err := s.boothSvc.Create(booth.id, booth.hall, booth.zone, booth.name, booth.watts); err != nil {
				return err
			}
		}
		if err := s.hallSvc.SetBoothPosition(booth.hall, booth.id, domain.Point{X: booth.posX, Y: booth.posY}); err != nil {
			return err
		}
	}
	if !s.quotaSvc.HasQuota("h-a") {
		if err := s.quotaSvc.SetQuota("q-a", "h-a", 5000); err != nil {
			return err
		}
	}
	if !s.quotaSvc.HasQuota("h-b") {
		if err := s.quotaSvc.SetQuota("q-b", "h-b", 6000); err != nil {
			return err
		}
	}
	gates := []struct {
		id   string
		zone string
		name string
	}{
		{"g-1", "z-a1", "一号门"},
		{"g-2", "z-a1", "二号门"},
		{"g-3", "z-b1", "三号门"},
	}
	for _, gate := range gates {
		if !s.gateSvc.HasGate(gate.id) {
			if err := s.gateSvc.Create(gate.id, gate.zone, gate.name); err != nil {
				return err
			}
		}
	}
	screens := []struct {
		id   string
		name string
	}{
		{"s-main", "主场馆大屏"},
		{"s-b", "B 馆大屏"},
	}
	for _, screen := range screens {
		if !s.screenSvc.HasScreen(screen.id) {
			if err := s.screenSvc.Create(screen.id, screen.name); err != nil {
				return err
			}
		}
		s.lightSvc.RegisterScreen(screen.id)
	}
	clips := []domain.Clip{
		{ID: "c-1", Name: "开场宣传", Sponsor: "冠名商", DurationSec: 30},
		{ID: "c-2", Name: "展商巡礼", Sponsor: "金牌赞助", DurationSec: 45},
		{ID: "c-3", Name: "闭馆提示", Sponsor: "组委会", DurationSec: 20},
	}
	if err := s.screenSvc.SetPlaylist("s-main", clips); err != nil {
		return err
	}
	if err := s.screenSvc.SetPlaylist("s-b", clips); err != nil {
		return err
	}
	scenes := []domain.Scene{
		{ZoneID: "z-a1", Name: "布展", Brightness: 70, Color: "warm"},
		{ZoneID: "z-a2", Name: "布展", Brightness: 70, Color: "warm"},
		{ZoneID: "z-b1", Name: "布展", Brightness: 70, Color: "warm"},
		{ZoneID: "z-b2", Name: "布展", Brightness: 70, Color: "warm"},
	}
	for _, scene := range scenes {
		if err := s.lightSvc.UpdateScene(scene, time.Now().Unix()); err != nil {
			return err
		}
	}
	cameras := []struct {
		id   string
		hall string
		name string
	}{
		{"cam-a", "h-a", "A 馆云台"},
		{"cam-b", "h-b", "B 馆云台"},
	}
	for _, camera := range cameras {
		if !s.cameraSvc.HasCamera(camera.id) {
			if err := s.cameraSvc.Create(camera.id, camera.hall, camera.name); err != nil {
				return err
			}
		}
	}
	if err := s.cameraSvc.SetPreset("cam-a", "b-101", 10, 20); err != nil {
		return err
	}
	if err := s.cameraSvc.SetPreset("cam-b", "b-201", 15, 25); err != nil {
		return err
	}
	if err := s.cameraSvc.RefreshPresets("cam-a", 0); err != nil {
		return err
	}
	if err := s.cameraSvc.RefreshPresets("cam-b", 0); err != nil {
		return err
	}
	routes := []struct {
		id     string
		zone   string
		name   string
	}{
		{"r-a", "z-a1", "A 馆巡更"},
		{"r-b", "z-b1", "B 馆巡更"},
	}
	for _, route := range routes {
		if !s.guardSvc.HasRoute(route.id) {
			if err := s.guardSvc.Create(route.id, route.zone, route.name); err != nil {
				return err
			}
		}
	}
	if err := s.store.SaveSeason(store.Season{Name: "autumn", NightHour: 18, NightMin: 30}); err != nil {
		return err
	}
	return nil
}

func (s *Server) SeedIfEmpty() error {
	count, err := s.hallSvc.CountHalls()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.Seed()
}
