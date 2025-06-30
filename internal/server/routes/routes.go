package routes

import (
	"net/http"

	"github.com/savin000/audit-log/internal/server/handlers"
)

func RegisterRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/audit-logs", h.GetAuditLogsHandler)
	return mux
}
