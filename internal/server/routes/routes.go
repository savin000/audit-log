package routes

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/savin000/audit-log/internal/server/handlers"
)

func RegisterRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/audit-logs", h.GetAuditLogsHandler)
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
