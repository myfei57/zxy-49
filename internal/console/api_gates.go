package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type gateRequest struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Name   string `json:"name"`
}

func (s *Server) handleListGates(w http.ResponseWriter, r *http.Request) {
	gates, err := s.gateSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gates)
}

func (s *Server) handleCreateGate(w http.ResponseWriter, r *http.Request) {
	var req gateRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.gateSvc.Create(req.ID, req.ZoneID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetGate(w http.ResponseWriter, r *http.Request) {
	gate, err := s.gateSvc.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gate)
}

func (s *Server) handleGateMode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Mode string `json:"mode"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.gateSvc.SetMode(chi.URLParam(r, "id"), req.Mode, time.Now().Unix()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": req.Mode})
}
