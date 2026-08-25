package console

import (
	"net/http"
	"strconv"

	"venueops/internal/domain"

	"github.com/go-chi/chi/v5"
)

type hallRequest struct {
	ID      string `json:"id"`
	VenueID string `json:"venue_id"`
	Name    string `json:"name"`
}

func (s *Server) handleListHalls(w http.ResponseWriter, r *http.Request) {
	halls, err := s.hallSvc.ListHalls()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, halls)
}

func (s *Server) handleCreateHall(w http.ResponseWriter, r *http.Request) {
	var req hallRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.CreateHall(req.ID, req.VenueID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetHall(w http.ResponseWriter, r *http.Request) {
	hall, err := s.hallSvc.GetHall(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, hall)
}

func (s *Server) handleLayoutChange(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Positions map[string]domain.Point `json:"positions"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.LayoutChange(chi.URLParam(r, "id"), req.Positions); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"layout": "updated"})
}

func (s *Server) handleBoothChange(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "id")
	boothID := chi.URLParam(r, "boothID")
	add := true
	if raw := r.URL.Query().Get("add"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		add = parsed
	}
	var req struct {
		ZoneID string `json:"zone_id"`
		Name   string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.hallSvc.BoothChange(hallID, boothID, req.ZoneID, req.Name, add); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"booth": boothID, "add": strconv.FormatBool(add)})
}
