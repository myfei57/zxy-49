package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type zoneRequest struct {
	ID      string `json:"id"`
	HallID  string `json:"hall_id"`
	VenueID string `json:"venue_id"`
	Name    string `json:"name"`
}

func (s *Server) handleListZones(w http.ResponseWriter, r *http.Request) {
	zones, err := s.zoneSvc.ListZones()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleCreateZone(w http.ResponseWriter, r *http.Request) {
	var req zoneRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.zoneSvc.CreateZone(req.ID, req.HallID, req.VenueID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.AddZone(req.HallID, req.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.nsSvc.AttachZone(req.VenueID, req.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetZone(w http.ResponseWriter, r *http.Request) {
	zone, err := s.zoneSvc.GetZone(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zone)
}

func (s *Server) handleSplitZone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HallID      string `json:"hall_id"`
		PartitionA  string `json:"partition_a"`
		PartitionB  string `json:"partition_b"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.zoneSvc.Split(req.HallID, req.PartitionA, req.PartitionB); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"split": "ok"})
}
