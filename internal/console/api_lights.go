package console

import (
	"net/http"
	"time"

	"venueops/internal/domain"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListScenes(w http.ResponseWriter, r *http.Request) {
	scenes, err := s.lightSvc.Scenes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, scenes)
}

func (s *Server) handleCreateScene(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ZoneID     string `json:"zone_id"`
		Name       string `json:"name"`
		Brightness int    `json:"brightness"`
		Color      string `json:"color"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	scene := domain.Scene{
		ZoneID:     req.ZoneID,
		Name:       req.Name,
		Brightness: req.Brightness,
		Color:      req.Color,
	}
	if err := s.lightSvc.UpdateScene(scene, time.Now().Unix()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, scene)
}

func (s *Server) handleEffectiveScene(w http.ResponseWriter, r *http.Request) {
	scene, err := s.lightSvc.EffectiveScene(chi.URLParam(r, "zoneID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, scene)
}

func (s *Server) handleEvacuate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HallID string `json:"hall_id"`
		Active bool   `json:"active"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now().Unix()
	var err error
	if req.Active {
		err = s.lightSvc.Evacuate(req.HallID, now)
	} else {
		err = s.lightSvc.ClearEvacuation(req.HallID, now)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"active": req.Active})
}

func (s *Server) handleEvacuationState(w http.ResponseWriter, r *http.Request) {
	state, err := s.lightSvc.EvacuationState(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, state)
}
