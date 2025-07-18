package api

import (
	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/middleware"
	"github.com/buiminhduc234/audit-log-api/internal/service"
	"github.com/buiminhduc234/audit-log-api/internal/service/pubsub"
	"github.com/buiminhduc234/audit-log-api/pkg/logger"
)

type Server struct {
	tenant    *TenantHandler
	auditLog  *AuditLogHandler
	websocket *WebSocketHandler
	auth      *middleware.AuthMiddleware
}

func NewServer(
	tenantService *service.TenantService,
	auditLogService *service.AuditLogService,
	auth *middleware.AuthMiddleware,
	logger *logger.Logger,
	pubsub *pubsub.RedisPubSub,
) *Server {
	return &Server{
		tenant:    NewTenantHandler(tenantService),
		auditLog:  NewAuditLogHandler(auditLogService),
		websocket: NewWebSocketHandler(auditLogService, logger, pubsub),
		auth:      auth,
	}
}

func (s *Server) SetupRoutes(api *gin.RouterGroup) {
	{
		tenants := api.Group("/tenants", s.auth.JWTAuth(), s.auth.RequireRole("admin"))
		{
			tenants.POST("", s.tenant.CreateTenant)
			tenants.GET("", s.tenant.ListTenants)
		}

		logs := api.Group("/logs", s.auth.JWTAuth(), s.auth.RequireRole("user"))
		{
			logs.POST("", s.auditLog.CreateLog)
			logs.GET("", s.auditLog.ListLogs)
			logs.GET("/:id", s.auditLog.GetLog)
			logs.GET("/export", s.auditLog.ExportLogs)
			logs.GET("/stats", s.auditLog.GetStats)
			logs.POST("/bulk", s.auditLog.BulkCreateLogs)
			logs.DELETE("/cleanup", s.auth.RequireRole("auditor"), s.auditLog.Cleanup)
			logs.GET("/stream", s.websocket.HandleWebSocket)
		}
	}
}

// StartWebSocketHub starts the WebSocket hub for broadcasting logs
func (s *Server) StartWebSocketHub() {
	go s.websocket.Start()
}

// GetWebSocketHandler returns the WebSocket handler for wiring up broadcasting
func (s *Server) GetWebSocketHandler() *WebSocketHandler {
	return s.websocket
}
