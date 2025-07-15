package api

import (
	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/middleware"
	"github.com/buiminhduc234/audit-log-api/internal/service"
)

type Server struct {
	tenant    *TenantHandler
	auditLog  *AuditLogHandler
	auth      *middleware.AuthMiddleware
	rateLimit *middleware.RateLimitMiddleware
}

func NewServer(
	tenantService *service.TenantService,
	auditLogService *service.AuditLogService,
	auth *middleware.AuthMiddleware,
	rateLimit *middleware.RateLimitMiddleware,
) *Server {
	return &Server{
		tenant:    NewTenantHandler(tenantService),
		auditLog:  NewAuditLogHandler(auditLogService),
		auth:      auth,
		rateLimit: rateLimit,
	}
}

func (s *Server) SetupRoutes(api *gin.RouterGroup) {
	{
		tenants := api.Group("/tenants", s.auth.JWTAuth(), s.auth.RequireRole("admin"))
		{
			tenants.POST("", s.tenant.CreateTenant)
			tenants.GET("", s.tenant.ListTenants)
		}

		logs := api.Group("/logs", s.auth.JWTAuth(), s.rateLimit.RateLimit(), s.auth.RequireRole("user"))
		{
			logs.POST("", s.auditLog.CreateLog)
			logs.GET("", s.auditLog.ListLogs)
			logs.GET("/:id", s.auditLog.GetLog)
			logs.GET("/export", s.auditLog.ExportLogs)
			logs.GET("/stats", s.auditLog.GetStats)
			logs.POST("/bulk", s.auditLog.CreateBulkLogs)
			logs.DELETE("/cleanup", s.auth.RequireRole("audit"), s.auditLog.CleanupLogs)
		}
	}
}
