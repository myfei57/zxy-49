package console

import (
	"net/http"
	"time"

	"venueops/internal/domain"
	"venueops/internal/store"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	venues, _ := s.nsSvc.CountVenues()
	halls, _ := s.hallSvc.CountHalls()
	zones, _ := s.zoneSvc.CountZones()
	booths, _ := s.boothSvc.BoothCount()
	gates, _ := s.gateSvc.Count()
	screens, _ := s.screenSvc.Count()
	cameras, _ := s.cameraSvc.Count()
	routes, _ := s.guardSvc.Count()
	quotas, _ := s.quotaSvc.ListQuotas()
	scenes, _ := s.lightSvc.SceneCount()
	auditSize, _ := s.auditSvc.Size()
	powered, _ := s.boothSvc.PoweredCount()
	totalLoad, _ := s.boothSvc.TotalLoad()
	poweredLoad, _ := s.boothSvc.PoweredLoad()
	writeJSON(w, http.StatusOK, map[string]int{
		"venues":       venues,
		"halls":        halls,
		"zones":        zones,
		"booths":       booths,
		"gates":        gates,
		"screens":      screens,
		"cameras":      cameras,
		"routes":       routes,
		"quotas":       len(quotas),
		"scenes":       scenes,
		"audit_size":   auditSize,
		"powered":      powered,
		"total_load":   totalLoad,
		"powered_load": poweredLoad,
	})
}

func (s *Server) handleVenueZones(w http.ResponseWriter, r *http.Request) {
	zones, err := s.nsSvc.ZonesOfVenue(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleVenueRename(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.nsSvc.RenameVenue(chi.URLParam(r, "id"), req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"renamed": req.Name})
}

func (s *Server) handleVenueEnsure(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.nsSvc.EnsureVenue(chi.URLParam(r, "id"), req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"ensured": chi.URLParam(r, "id")})
}

func (s *Server) handleVenueValidate(w http.ResponseWriter, r *http.Request) {
	if err := s.nsSvc.ValidateVenue(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

func (s *Server) handleVenueNames(w http.ResponseWriter, r *http.Request) {
	names, err := s.nsSvc.VenueNames()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, names)
}

func (s *Server) handleVenueCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.nsSvc.CountVenues()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleHallsOfVenue(w http.ResponseWriter, r *http.Request) {
	halls, err := s.hallSvc.HallsOfVenue(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, halls)
}

func (s *Server) handleHallRename(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.RenameHall(chi.URLParam(r, "id"), req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"renamed": req.Name})
}

func (s *Server) handleHallAttachZone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZoneID string `json:"zone_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.AttachVenueZone(chi.URLParam(r, "id"), req.ZoneID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"attached": req.ZoneID})
}

func (s *Server) handleHallZoneIDs(w http.ResponseWriter, r *http.Request) {
	zones, err := s.hallSvc.HallZoneIDs(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleHallRemoveZone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZoneID string `json:"zone_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.RemoveZone(chi.URLParam(r, "id"), req.ZoneID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"removed": req.ZoneID})
}

func (s *Server) handleHallLayoutSize(w http.ResponseWriter, r *http.Request) {
	size, err := s.hallSvc.LayoutSize(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"size": size})
}

func (s *Server) handleHallMoveBooth(w http.ResponseWriter, r *http.Request) {
	var req domain.Point
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.MoveBooth(chi.URLParam(r, "id"), chi.URLParam(r, "boothID"), req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"moved": chi.URLParam(r, "boothID")})
}

func (s *Server) handleHallRemoveBoothPosition(w http.ResponseWriter, r *http.Request) {
	if err := s.hallSvc.RemoveBoothPosition(chi.URLParam(r, "id"), chi.URLParam(r, "boothID")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"removed": chi.URLParam(r, "boothID")})
}

func (s *Server) handleHallBoothExists(w http.ResponseWriter, r *http.Request) {
	exists := s.hallSvc.BoothExists(chi.URLParam(r, "id"), chi.URLParam(r, "boothID"))
	writeJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

func (s *Server) handleHallBoothZones(w http.ResponseWriter, r *http.Request) {
	zones, err := s.hallSvc.BoothZones(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleHallZoneCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.hallSvc.ZoneCount(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleHallPowerOnAll(w http.ResponseWriter, r *http.Request) {
	on, err := s.boothSvc.PowerOnAll(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"powered_on": on})
}

func (s *Server) handleHallPowerOffAll(w http.ResponseWriter, r *http.Request) {
	off, err := s.boothSvc.PowerOffAll(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"powered_off": off})
}

func (s *Server) handleHallLoad(w http.ResponseWriter, r *http.Request) {
	load, err := s.boothSvc.HallLoad(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"load": load})
}

func (s *Server) handleZoneVenue(w http.ResponseWriter, r *http.Request) {
	venueID, err := s.zoneSvc.VenueOfZone(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"venue_id": venueID})
}

func (s *Server) handleZoneRename(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.zoneSvc.RenameZone(chi.URLParam(r, "id"), req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"renamed": req.Name})
}

func (s *Server) handleZoneSetPartition(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Partition string `json:"partition"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.zoneSvc.SetPartition(chi.URLParam(r, "id"), req.Partition); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"partition": req.Partition})
}

func (s *Server) handleZoneDim(w http.ResponseWriter, r *http.Request) {
	scene, err := s.zoneSvc.Dim(chi.URLParam(r, "id"), s.lightSvc.SceneFor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, scene)
}

func (s *Server) handleZoneDimLevel(w http.ResponseWriter, r *http.Request) {
	scene, err := s.zoneSvc.DimLevel(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, scene)
}

func (s *Server) handleZoneSetDimLevel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Brightness int `json:"brightness"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.zoneSvc.SetDimLevel(chi.URLParam(r, "id"), req.Brightness); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"brightness": req.Brightness})
}

func (s *Server) handleZoneClearDim(w http.ResponseWriter, r *http.Request) {
	if err := s.zoneSvc.ClearDim(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"cleared": "ok"})
}

func (s *Server) handleHallDimAll(w http.ResponseWriter, r *http.Request) {
	results, err := s.zoneSvc.DimAll(chi.URLParam(r, "id"), s.lightSvc.SceneFor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleHallPartitions(w http.ResponseWriter, r *http.Request) {
	partitions, err := s.zoneSvc.PartitionsOfHall(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, partitions)
}

func (s *Server) handleHallZonesInPartition(w http.ResponseWriter, r *http.Request) {
	zones, err := s.zoneSvc.ZonesInPartition(chi.URLParam(r, "id"), chi.URLParam(r, "partition"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleHallPartitionCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.zoneSvc.PartitionCount(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleLightSceneNames(w http.ResponseWriter, r *http.Request) {
	names, err := s.lightSvc.SceneNames()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, names)
}

func (s *Server) handleLightSceneCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.lightSvc.SceneCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleLightScheduled(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZoneID string `json:"zone_id"`
		At     int64  `json:"at"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.lightSvc.ApplyScheduled(req.ZoneID, req.At); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"scheduled": req.ZoneID})
}

func (s *Server) handleLightEvacuating(w http.ResponseWriter, r *http.Request) {
	halls, err := s.lightSvc.EvacuatingHalls()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, halls)
}

func (s *Server) handleLightHandoverAll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Emergency bool  `json:"emergency"`
		At        int64 `json:"at"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.lightSvc.HandoverAll(req.Emergency, req.At); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"emergency": req.Emergency})
}

func (s *Server) handleLightScreenCount(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]int{"count": s.lightSvc.ScreenCount()})
}

func (s *Server) handleLightScreens(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.lightSvc.Screens())
}

func (s *Server) handleScreenPlaylistSize(w http.ResponseWriter, r *http.Request) {
	size, err := s.screenSvc.PlaylistSize(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"size": size})
}

func (s *Server) handleScreenSponsors(w http.ResponseWriter, r *http.Request) {
	sponsors, err := s.screenSvc.Sponsors(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sponsors)
}

func (s *Server) handleScreenHandover(w http.ResponseWriter, r *http.Request) {
	emergency, err := s.screenSvc.HandoverState(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"emergency": emergency})
}

func (s *Server) handleScreensHandover(w http.ResponseWriter, r *http.Request) {
	screens, err := s.screenSvc.HandoverScreens()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, screens)
}

func (s *Server) handleScreenClearHandover(w http.ResponseWriter, r *http.Request) {
	if err := s.screenSvc.ClearHandover(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"cleared": "ok"})
}

func (s *Server) handleScreenResetAt(w http.ResponseWriter, r *http.Request) {
	reset, err := s.screenSvc.ResetAt(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reset)
}

func (s *Server) handleScreenResetCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.screenSvc.ResetCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleScreenIDs(w http.ResponseWriter, r *http.Request) {
	ids, err := s.screenSvc.ScreenIDs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (s *Server) handleScreenNames(w http.ResponseWriter, r *http.Request) {
	names, err := s.screenSvc.Names()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, names)
}

func (s *Server) handleScreenPlaylistExists(w http.ResponseWriter, r *http.Request) {
	exists := s.store.ClipsExist(chi.URLParam(r, "id"))
	writeJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

func (s *Server) handleBoothSetLoad(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Watts int `json:"watts"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.boothSvc.SetLoad(chi.URLParam(r, "id"), req.Watts); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"watts": req.Watts})
}

func (s *Server) handleBoothToggle(w http.ResponseWriter, r *http.Request) {
	powered, err := s.boothSvc.PowerToggle(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"powered": powered})
}

func (s *Server) handleBoothPoweredCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.boothSvc.PoweredCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleBoothTotalLoad(w http.ResponseWriter, r *http.Request) {
	load, err := s.boothSvc.TotalLoad()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"load": load})
}

func (s *Server) handleBoothPoweredLoad(w http.ResponseWriter, r *http.Request) {
	load, err := s.boothSvc.PoweredLoad()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"load": load})
}

func (s *Server) handleGateArmed(w http.ResponseWriter, r *http.Request) {
	ids, err := s.gateSvc.ArmedGates()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (s *Server) handleGateReleasedList(w http.ResponseWriter, r *http.Request) {
	ids, err := s.gateSvc.ReleasedGates()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (s *Server) handleGateState(w http.ResponseWriter, r *http.Request) {
	state, err := s.gateSvc.ArmingState(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": state})
}

func (s *Server) handleGateDurable(w http.ResponseWriter, r *http.Request) {
	mode, err := s.gateSvc.DurableMode(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": mode})
}

func (s *Server) handleGateDisarm(w http.ResponseWriter, r *http.Request) {
	if err := s.gateSvc.Disarm(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"disarmed": "ok"})
}

func (s *Server) handleGateSetModeByName(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Mode string `json:"mode"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.gateSvc.SetModeByName(req.Name, req.Mode, time.Now().Unix()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": req.Mode})
}

func (s *Server) handleGateZone(w http.ResponseWriter, r *http.Request) {
	zoneID, err := s.gateSvc.ZoneOf(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"zone_id": zoneID})
}

func (s *Server) handleGatesOfZone(w http.ResponseWriter, r *http.Request) {
	gates, err := s.gateSvc.GatesOfZone(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gates)
}

func (s *Server) handleGateJournalExists(w http.ResponseWriter, r *http.Request) {
	exists := s.store.GateModeExists(chi.URLParam(r, "id"))
	writeJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

func (s *Server) handleCameraPresetSingle(w http.ResponseWriter, r *http.Request) {
	preset, err := s.cameraSvc.PresetFor(chi.URLParam(r, "id"), chi.URLParam(r, "boothID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, preset)
}

func (s *Server) handleCameraPresetCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.cameraSvc.PresetCount(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleCameraPatrolBooths(w http.ResponseWriter, r *http.Request) {
	booths, err := s.cameraSvc.PatrolBooths(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, booths)
}

func (s *Server) handleRouteCheckpointCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.guardSvc.CheckpointCount(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleRouteRebuildNew(w http.ResponseWriter, r *http.Request) {
	var req routeRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.guardSvc.RebuildRoute(req.ID, req.ZoneID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleSeasonExists(w http.ResponseWriter, r *http.Request) {
	exists := s.store.SeasonExists()
	writeJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

func (s *Server) handleAuditAll(w http.ResponseWriter, r *http.Request) {
	records, err := s.auditSvc.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditSize(w http.ResponseWriter, r *http.Request) {
	size, err := s.auditSvc.Size()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"size": size})
}

func (s *Server) handleAuditCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.auditSvc.Count(r.URL.Query().Get("action"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleZoneDimStatus(w http.ResponseWriter, r *http.Request) {
	zones, err := s.zoneSvc.CountZones()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"zones": zones})
}

func (s *Server) handleBoothExists(w http.ResponseWriter, r *http.Request) {
	exists := s.boothSvc.HasBooth(chi.URLParam(r, "id"))
	writeJSON(w, http.StatusOK, map[string]bool{"exists": exists})
}

func (s *Server) handleScreenPosition(w http.ResponseWriter, r *http.Request) {
	position, err := s.screenSvc.Position(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"position": position})
}

func (s *Server) handleScreenClips(w http.ResponseWriter, r *http.Request) {
	clips, err := s.screenSvc.Clips(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, clips)
}

func (s *Server) handleCameraSeason(w http.ResponseWriter, r *http.Request) {
	season, err := s.store.LoadSeason()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, season)
}

func (s *Server) handleCameraSetSeason(w http.ResponseWriter, r *http.Request) {
	var req store.Season
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.SaveSeason(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"season": req.Name})
}
