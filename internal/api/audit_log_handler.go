package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/service"
)

type AuditLogHandler struct {
	*BaseHandler
	service *service.AuditLogService
}

func NewAuditLogHandler(service *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

// CreateLog godoc
// @Summary Create a new audit log
// @Description Create a new audit log entry
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param log body domain.AuditLog true "Audit log object"
// @Success 201 {object} domain.AuditLog
// @Failure 400 {object} ErrorResponse
// @Router /logs [post]
func (h *AuditLogHandler) CreateLog(c *gin.Context) {
	var log domain.AuditLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(h.RequestCtx(c), &log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// BulkCreateLogs godoc
// @Summary Create multiple audit logs
// @Description Create multiple audit log entries in a single request
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param logs body []domain.AuditLog true "Array of audit log objects"
// @Success 201 {array} domain.AuditLog
// @Failure 400 {object} ErrorResponse
// @Router /logs/bulk [post]
func (h *AuditLogHandler) BulkCreateLogs(c *gin.Context) {
	var logs []domain.AuditLog
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BulkCreate(h.RequestCtx(c), logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, logs)
}

// GetLog godoc
// @Summary Get a specific audit log
// @Description Get an audit log by ID
// @Tags audit-logs
// @Produce json
// @Param id path string true "Log ID"
// @Success 200 {object} domain.AuditLog
// @Failure 404 {object} ErrorResponse
// @Router /logs/{id} [get]
func (h *AuditLogHandler) GetLog(c *gin.Context) {
	id := c.Param("id")

	log, err := h.service.GetByID(h.RequestCtx(c), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListLogs godoc
// @Summary List audit logs
// @Description Get a list of audit logs with filtering options
// @Tags audit-logs
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param user_id query string false "Filter by user ID"
// @Param action query string false "Filter by action"
// @Param resource query string false "Filter by resource"
// @Param severity query string false "Filter by severity"
// @Param start_time query string false "Filter by start time (RFC3339)"
// @Param end_time query string false "Filter by end time (RFC3339)"
// @Success 200 {array} domain.AuditLog
// @Router /logs [get]
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	filter := &domain.AuditLogFilter{
		UserID:       c.Query("user_id"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		Severity:     c.Query("severity"),
	}

	// Parse pagination
	if page := c.Query("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			filter.Page = pageNum
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			filter.PageSize = size
		}
	}

	// Parse time filters
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filter.EndTime = t
		}
	}

	logs, err := h.service.List(h.RequestCtx(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// DeleteLog godoc
// @Summary Delete an audit log
// @Description Delete an audit log by ID
// @Tags audit-logs
// @Param id path string true "Log ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /logs/{id} [delete]
func (h *AuditLogHandler) DeleteLog(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(h.RequestCtx(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteOldLogs godoc
// @Summary Delete old audit logs
// @Description Delete audit logs older than specified days
// @Tags audit-logs
// @Param days query int true "Number of days to retain"
// @Success 204 "No Content"
// @Router /logs/cleanup [delete]
func (h *AuditLogHandler) DeleteOldLogs(c *gin.Context) {
	days := 90 // Default to 90 days

	if err := h.service.DeleteOlderThan(h.RequestCtx(c), days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportLogs godoc
// @Summary Export audit logs
// @Description Export audit logs in JSON or CSV format
// @Tags audit-logs
// @Produce json,text/csv
// @Param format query string false "Export format (json or csv)" default(json)
// @Param filter query domain.AuditLogFilter false "Filter parameters"
// @Success 200 {file} file
// @Failure 400 {object} ErrorResponse
// @Router /logs/export [get]
func (h *AuditLogHandler) ExportLogs(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Must be 'json' or 'csv'"})
		return
	}

	// Get tenant ID from context
	tenantID, _ := c.Get("tenant_id")
	filter := &domain.AuditLogFilter{
		TenantID:     tenantID.(string),
		UserID:       c.Query("user_id"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		Severity:     c.Query("severity"),
	}

	// Parse time filters
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filter.EndTime = t
		}
	}

	logs, err := h.service.List(h.RequestCtx(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch format {
	case "json":
		c.Header("Content-Disposition", "attachment; filename=audit_logs.json")
		c.JSON(http.StatusOK, logs)
	case "csv":
		c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
		c.Header("Content-Type", "text/csv")

		// Write CSV header
		c.Writer.Write([]byte("ID,TenantID,UserID,SessionID,Action,Resource,ResourceID,IPAddress,UserAgent,Severity,Timestamp\n"))

		// Write each log entry as CSV
		for _, log := range logs {
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
				log.ID, log.TenantID, log.UserID, log.SessionID, log.Action,
				log.ResourceType, log.ResourceID, log.IPAddress, log.UserAgent,
				log.Severity, log.Timestamp.Format(time.RFC3339))
			c.Writer.Write([]byte(row))
		}
	}
}

// GetStats godoc
// @Summary Get audit log statistics
// @Description Get statistics about audit logs including counts by action, severity, and resource
// @Tags audit-logs
// @Produce json
// @Param start_time query string false "Filter by start time (RFC3339)"
// @Param end_time query string false "Filter by end time (RFC3339)"
// @Success 200 {object} domain.AuditLogStats
// @Router /logs/stats [get]
func (h *AuditLogHandler) GetStats(c *gin.Context) {
	filter := &domain.AuditLogFilter{}

	// Parse time filters
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filter.EndTime = t
		}
	}

	stats, err := h.service.GetStats(h.RequestCtx(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CreateBulkLogs godoc
// @Summary Create multiple audit logs
// @Description Create multiple audit log entries in a single request
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param logs body []domain.AuditLog true "Array of audit log objects"
// @Success 201 {array} domain.AuditLog
// @Failure 400 {object} ErrorResponse
// @Router /logs/bulk [post]
func (h *AuditLogHandler) CreateBulkLogs(c *gin.Context) {
	var logs []domain.AuditLog
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BulkCreate(h.RequestCtx(c), logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, logs)
}

// CleanupLogs godoc
// @Summary Delete old audit logs
// @Description Delete audit logs older than specified retention period
// @Tags audit-logs
// @Param days query int false "Number of days to retain" default(90)
// @Success 204 "No Content"
// @Router /logs/cleanup [delete]
func (h *AuditLogHandler) CleanupLogs(c *gin.Context) {
	days := 90 // Default retention period

	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	if err := h.service.DeleteOlderThan(h.RequestCtx(c), days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
