package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListCameras(w http.ResponseWriter, r *http.Request) {
	cameras, err := s.cameraSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cameras)
}

func (s *Server) handleCreateCamera(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string `json:"id"`
		HallID string `json:"hall_id"`
		Name   string `json:"name"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cameraSvc.Create(req.ID, req.HallID, req.Name); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetCamera(w http.ResponseWriter, r *http.Request) {
	camera, err := s.cameraSvc.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, camera)
}

func (s *Server) handleCameraPreset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BoothID string `json:"booth_id"`
		Pan     int    `json:"pan"`
		Tilt    int    `json:"tilt"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cameraSvc.SetPreset(chi.URLParam(r, "id"), req.BoothID, req.Pan, req.Tilt); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"preset": req.BoothID})
}

func (s *Server) handleCameraRefresh(w http.ResponseWriter, r *http.Request) {
	if err := s.cameraSvc.RefreshPresets(chi.URLParam(r, "id"), 0); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"refresh": "ok"})
}

func (s *Server) handleCameraDayNight(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	mode, err := s.cameraSvc.DayNightSwitch(chi.URLParam(r, "id"), req.Hour, req.Minute)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": mode})
}

func (s *Server) handleCameraPresets(w http.ResponseWriter, r *http.Request) {
	booths, err := s.cameraSvc.PresetBooths(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, booths)
}

func (s *Server) handleCameraPresetPosition(w http.ResponseWriter, r *http.Request) {
	position, err := s.cameraSvc.PresetPosition(chi.URLParam(r, "id"), chi.URLParam(r, "boothID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, position)
}

func (s *Server) handleCameraRemovePreset(w http.ResponseWriter, r *http.Request) {
	if err := s.cameraSvc.RemovePreset(chi.URLParam(r, "id"), chi.URLParam(r, "boothID")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"removed": chi.URLParam(r, "boothID")})
}

func (s *Server) handleCameraDayNightState(w http.ResponseWriter, r *http.Request) {
	mode, err := s.cameraSvc.DayNightState(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": mode})
}

func (s *Server) handleCameraSetDayNight(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Mode string `json:"mode"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cameraSvc.SetDayNight(chi.URLParam(r, "id"), req.Mode); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mode": req.Mode})
}

func (s *Server) handleCameraPatrol(w http.ResponseWriter, r *http.Request) {
	var req struct {
		At int64 `json:"at"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	preset, err := s.cameraSvc.Patrol(chi.URLParam(r, "id"), req.At)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, preset)
}

func (s *Server) handleCameraPatrolTarget(w http.ResponseWriter, r *http.Request) {
	at := int64(0)
	if raw := r.URL.Query().Get("at"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		at = parsed
	}
	target, err := s.cameraSvc.PatrolTarget(chi.URLParam(r, "id"), at)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"target": target})
}

func (s *Server) handleCameraPatrolRounds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Count int `json:"count"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rounds, err := s.cameraSvc.PatrolRounds(chi.URLParam(r, "id"), req.Count)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rounds)
}

func (s *Server) handleCameraPatrolSequence(w http.ResponseWriter, r *http.Request) {
	sequence, err := s.cameraSvc.PatrolSequence(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sequence)
}

func (s *Server) handleCameraNight(w http.ResponseWriter, r *http.Request) {
	ids, err := s.cameraSvc.NightCameras()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (s *Server) handleCameraDay(w http.ResponseWriter, r *http.Request) {
	ids, err := s.cameraSvc.DayCameras()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (s *Server) handleCamerasOfHall(w http.ResponseWriter, r *http.Request) {
	cameras, err := s.cameraSvc.CamerasOfHall(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cameras)
}

func (s *Server) handleCameraHall(w http.ResponseWriter, r *http.Request) {
	hallID, err := s.cameraSvc.HallOf(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"hall_id": hallID})
}
