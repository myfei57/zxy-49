package console

import (
	"net/http"
	"strconv"

	"venueops/internal/domain"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleGetQuota(w http.ResponseWriter, r *http.Request) {
	quota, err := s.quotaSvc.GetQuota(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, quota)
}

func (s *Server) handleSetQuota(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID         string `json:"id"`
		HallID     string `json:"hall_id"`
		LimitWatts int    `json:"limit_watts"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.quotaSvc.SetQuota(req.ID, req.HallID, req.LimitWatts); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": req.ID})
}

func (s *Server) handleListQuotas(w http.ResponseWriter, r *http.Request) {
	quotas, err := s.quotaSvc.ListQuotas()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, quotas)
}

func (s *Server) handleQuotaLimit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LimitWatts int `json:"limit_watts"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.quotaSvc.ChangeLimit(chi.URLParam(r, "hallID"), req.LimitWatts); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"limit": "updated"})
}

func (s *Server) handleQuotaReset(w http.ResponseWriter, r *http.Request) {
	if err := s.quotaSvc.ResetUsage(chi.URLParam(r, "hallID")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"reset": "ok"})
}

func (s *Server) handleQuotaAvailable(w http.ResponseWriter, r *http.Request) {
	available, err := s.quotaSvc.Available(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"available": available})
}

func (s *Server) handleQuotaCanReserve(w http.ResponseWriter, r *http.Request) {
	watts := 0
	if raw := r.URL.Query().Get("watts"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		watts = parsed
	}
	can, err := s.quotaSvc.CanReserve(chi.URLParam(r, "hallID"), watts)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"can_reserve": can})
}

func (s *Server) handleQuotaUsage(w http.ResponseWriter, r *http.Request) {
	used, limit, err := s.quotaSvc.Usage(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"used": used, "limit": limit})
}

func (s *Server) handleQuotaLedger(w http.ResponseWriter, r *http.Request) {
	ledger, err := s.quotaSvc.Ledger(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ledger)
}

func (s *Server) handleQuotaAppendLedger(w http.ResponseWriter, r *http.Request) {
	var quota domain.Quota
	if err := decodeBody(r, &quota); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.quotaSvc.AppendLedger(chi.URLParam(r, "hallID"), quota); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"appended": quota.ID})
}

func (s *Server) handleQuotaLedgerSize(w http.ResponseWriter, r *http.Request) {
	size, err := s.quotaSvc.LedgerSize(chi.URLParam(r, "hallID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"size": size})
}
