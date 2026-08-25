package console

import (
	"net/http"
	"time"

	"venueops/internal/domain"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListScreens(w http.ResponseWriter, r *http.Request) {
	screens, err := s.screenSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, screens)
}

func (s *Server) handleCreateScreen(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.screenSvc.Create(req.ID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.lightSvc.RegisterScreen(req.ID)
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetScreen(w http.ResponseWriter, r *http.Request) {
	screen, err := s.screenSvc.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, screen)
}

func (s *Server) handleScreenCurrent(w http.ResponseWriter, r *http.Request) {
	clip, err := s.screenSvc.Current(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, clip)
}

func (s *Server) handleScreenPlaylist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Clips []domain.Clip `json:"clips"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.screenSvc.SetPlaylist(chi.URLParam(r, "id"), req.Clips); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"playlist": "updated"})
}

func (s *Server) handleScreenTick(w http.ResponseWriter, r *http.Request) {
	clip, err := s.screenSvc.Tick(chi.URLParam(r, "id"), time.Now())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, clip)
}

func (s *Server) handleScreenReset(w http.ResponseWriter, r *http.Request) {
	if err := s.screenSvc.Reset(chi.URLParam(r, "id"), time.Now().Unix()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"reset": "ok"})
}
