package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/savin000/audit-log/internal/clickhouse"
	"net/http"
	"strconv"
)

type Handler struct {
	Ch *clickhouse.Client
}

func (h *Handler) GetAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 1 {
		offset = 0
	}

	auditLogs, err := h.Ch.GetAuditLogs(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get audit logs - %v", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(auditLogs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response - %v", err.Error()), http.StatusInternalServerError)
	}
}
