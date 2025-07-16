package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/api/dto"
	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/service"
	"github.com/buiminhduc234/audit-log-api/pkg/utils"
)

type AuditLogHandler struct {
	*BaseHandler
	service *service.AuditLogService
}

func NewAuditLogHandler(service *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

// CreateLog Create a new audit log entry
// @Summary Create audit log
// @Description Create a new audit log entry
// @Tags    audit_logs
// @Accept  json
// @Produce json
// @Param   body body dto.CreateAuditLogRequest true "Audit log object"
// @Success 201 {object}
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs [post]
func (h *AuditLogHandler) CreateLog(c *gin.Context) {
	var log dto.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	if err := h.service.Create(h.RequestCtx(c), log); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Log created successfully"})
}

// BulkCreateLogs Create multiple audit log entries
// @Summary Bulk create audit logs
// @Description Create multiple audit log entries in a single request
// @Tags    audit_logs
// @Accept  json
// @Produce json
// @Param   body body []dto.CreateAuditLogRequest true "Array of audit log objects"
// @Success 201 {object}
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs/bulk [post]
func (h *AuditLogHandler) BulkCreateLogs(c *gin.Context) {
	var logs []dto.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	if err := h.service.BulkCreate(h.RequestCtx(c), logs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Logs created successfully"})
}

// GetLog Get a specific audit log by ID
// @Summary Get audit log
// @Description Get an audit log entry by its ID
// @Tags    audit_logs
// @Produce json
// @Param   id path string true "Log ID"
// @Success 200 {object} dto.AuditLogResponse
// @Failure 401 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs/{id} [get]
func (h *AuditLogHandler) GetLog(c *gin.Context) {
	id := c.Param("id")

	log, err := h.service.GetByID(h.RequestCtx(c), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}
	if log == nil {
		c.JSON(http.StatusNotFound, dto.Error{Error: "Log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListLogs Get a list of audit logs with filtering
// @Summary List audit logs
// @Description Get a list of audit logs with filtering options
// @Tags    audit_logs
// @Produce json
// @Param   page query int false "Page number"
// @Param   page_size query int false "Page size"
// @Param   user_id query string false "Filter by user ID"
// @Param   action query string false "Filter by action"
// @Param   resource_type query string false "Filter by resource type"
// @Param   severity query string false "Filter by severity"
// @Param   start_time query string false "Filter by start time (RFC3339)"
// @Param   end_time query string false "Filter by end time (RFC3339)"
// @Success 200 {array} dto.AuditLogResponse
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs [get]
func (h *AuditLogHandler) ListLogs(c *gin.Context) {
	filter := &domain.AuditLogFilter{
		UserID:       c.Query("user_id"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		Severity:     c.Query("severity"),
		SessionID:    c.Query("session_id"),
		IPAddress:    c.Query("ip_address"),
		UserAgent:    c.Query("user_agent"),
		Message:      c.Query("message"),
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
		t, err := utils.ParseUserTime(startTime, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error{
				Error: err.Error(),
			})
			return
		}
		filter.StartTime = t
	}
	if endTime := c.Query("end_time"); endTime != "" {
		t, err := utils.ParseUserTime(endTime, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error{
				Error: err.Error(),
			})
			return
		}
		filter.EndTime = t
	}

	logs, err := h.service.List(h.RequestCtx(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// DeleteOldLogs Delete audit logs older than retention period
// @Summary Delete old logs
// @Description Delete audit logs older than specified days
// @Tags    audit_logs
// @Produce json
// @Param   days query int true "Number of days to retain"
// @Success 204 "No Content"
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs/cleanup [delete]
func (h *AuditLogHandler) DeleteOldLogs(c *gin.Context) {
	days := 90 // Default to 90 days

	if err := h.service.DeleteOlderThan(h.RequestCtx(c), days); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportLogs Export audit logs in JSON or CSV format
// @Summary Export audit logs
// @Description Export audit logs with filtering options in JSON or CSV format
// @Tags    audit_logs
// @Produce json,text/csv
// @Param   format query string false "Export format (json or csv)" default(json)
// @Param   user_id query string false "Filter by user ID"
// @Param   action query string false "Filter by action"
// @Param   resource_type query string false "Filter by resource type"
// @Param   severity query string false "Filter by severity"
// @Param   start_time query string false "Filter by start time (RFC3339)"
// @Param   end_time query string false "Filter by end time (RFC3339)"
// @Success 200 {file} file
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs/export [get]
func (h *AuditLogHandler) ExportLogs(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "Invalid format. Must be 'json' or 'csv'"})
		return
	}

	filter := &domain.AuditLogFilter{
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
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
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

// GetStats Get audit log statistics
// @Summary Get log statistics
// @Description Get statistics about audit logs including counts by action, severity, and resource
// @Tags    audit_logs
// @Produce json
// @Param   start_time query string false "Filter by start time (RFC3339)"
// @Param   end_time query string false "Filter by end time (RFC3339)"
// @Success 200 {object} dto.GetAuditLogStatsResponse
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router  /logs/stats [get]
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
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
