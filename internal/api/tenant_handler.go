package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/service"
)

type TenantHandler struct {
	*BaseHandler
	service *service.TenantService
}

func NewTenantHandler(service *service.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body domain.Tenant true "Tenant object"
// @Success 201 {object} domain.Tenant
// @Router /tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var tenant domain.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(h.RequestCtx(c), &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get a list of all tenants
// @Tags tenants
// @Produce json
// @Success 200 {array} domain.Tenant
// @Router /tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	tenants, err := h.service.List(h.RequestCtx(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenants)
}
