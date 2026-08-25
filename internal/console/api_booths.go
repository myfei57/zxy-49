package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type boothRequest struct {
	ID        string `json:"id"`
	HallID    string `json:"hall_id"`
	ZoneID    string `json:"zone_id"`
	Name      string `json:"name"`
	LoadWatts int    `json:"load_watts"`
}

func (s *Server) handleListBooths(w http.ResponseWriter, r *http.Request) {
	booths, err := s.boothSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, booths)
}

func (s *Server) handleCreateBooth(w http.ResponseWriter, r *http.Request) {
	var req boothRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.boothSvc.Create(req.ID, req.HallID, req.ZoneID, req.Name, req.LoadWatts); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleGetBooth(w http.ResponseWriter, r *http.Request) {
	booth, err := s.boothSvc.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, booth)
}

func (s *Server) handleBoothPower(w http.ResponseWriter, r *http.Request) {
	var req struct {
		On bool `json:"on"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var err error
	if req.On {
		err = s.boothSvc.PowerOn(chi.URLParam(r, "id"))
	} else {
		err = s.boothSvc.PowerOff(chi.URLParam(r, "id"))
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"powered": strconv.FormatBool(req.On)})
}
