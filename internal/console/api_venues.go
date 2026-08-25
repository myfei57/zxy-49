package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type venueRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) handleListVenues(w http.ResponseWriter, r *http.Request) {
	venues, err := s.nsSvc.ListVenues()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, venues)
}

func (s *Server) handleCreateVenue(w http.ResponseWriter, r *http.Request) {
	var req venueRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.nsSvc.CreateVenue(req.ID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetVenue(w http.ResponseWriter, r *http.Request) {
	venue, err := s.nsSvc.GetVenue(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, venue)
}

func (s *Server) handleAttachZone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZoneID string `json:"zone_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.nsSvc.AttachZone(chi.URLParam(r, "id"), req.ZoneID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"attached": req.ZoneID})
}
