package console

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type routeRequest struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Name   string `json:"name"`
}

func (s *Server) handleListRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := s.guardSvc.ListRoutes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, routes)
}

func (s *Server) handleCreateRoute(w http.ResponseWriter, r *http.Request) {
	var req routeRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.guardSvc.Create(req.ID, req.ZoneID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetRoute(w http.ResponseWriter, r *http.Request) {
	route, err := s.guardSvc.GetRoute(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, route)
}

func (s *Server) handleRebuildRoute(w http.ResponseWriter, r *http.Request) {
	if err := s.guardSvc.RebuildExisting(chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"rebuild": "ok"})
}

func (s *Server) handleCheckIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BoothID string `json:"booth_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.guardSvc.CheckIn(chi.URLParam(r, "id"), req.BoothID, time.Now().Unix()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"checked": req.BoothID})
}

func (s *Server) handleCheckedIn(w http.ResponseWriter, r *http.Request) {
	checked, err := s.guardSvc.CheckedIn(chi.URLParam(r, "id"), chi.URLParam(r, "boothID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"checked": checked})
}

func (s *Server) handleRouteProgress(w http.ResponseWriter, r *http.Request) {
	visited, total, err := s.guardSvc.Progress(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"visited": visited, "total": total})
}

func (s *Server) handleRouteOverdue(w http.ResponseWriter, r *http.Request) {
	limitSec := int64(300)
	if raw := r.URL.Query().Get("limit_sec"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		limitSec = parsed
	}
	overdue, err := s.guardSvc.Overdue(chi.URLParam(r, "id"), time.Now().Unix(), limitSec)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"overdue": overdue})
}

func (s *Server) handleRouteAlarm(w http.ResponseWriter, r *http.Request) {
	target, err := s.guardSvc.RouteAlarm(chi.URLParam(r, "zoneID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"target": target})
}

func (s *Server) handleRoutesOfZone(w http.ResponseWriter, r *http.Request) {
	routes, err := s.guardSvc.RoutesOfZone(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, routes)
}
