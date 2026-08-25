package console

import (
	"net/http"
	"strconv"
)

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	records, err := s.auditSvc.Recent(100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditLast(w http.ResponseWriter, r *http.Request) {
	record, err := s.auditSvc.Last()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleAuditActions(w http.ResponseWriter, r *http.Request) {
	actions, err := s.auditSvc.Actions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actions)
}

func (s *Server) handleAuditTargets(w http.ResponseWriter, r *http.Request) {
	targets, err := s.auditSvc.DistinctTargets()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, targets)
}

func (s *Server) handleAuditByAction(w http.ResponseWriter, r *http.Request) {
	records, err := s.auditSvc.ByAction(r.URL.Query().Get("action"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditByTarget(w http.ResponseWriter, r *http.Request) {
	records, err := s.auditSvc.ByTarget(r.URL.Query().Get("target"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditSince(w http.ResponseWriter, r *http.Request) {
	at := int64(0)
	if raw := r.URL.Query().Get("at"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		at = parsed
	}
	records, err := s.auditSvc.Since(at)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditRecord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
		Target string `json:"target"`
		Detail string `json:"detail"`
	}
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.auditSvc.RecordNow(req.Action, req.Target, req.Detail); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"recorded": req.Action})
}

func (s *Server) handleAuditClear(w http.ResponseWriter, r *http.Request) {
	if err := s.auditSvc.Clear(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"cleared": "ok"})
}
