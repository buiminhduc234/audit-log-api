package utils

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey represents a key in the gin context
type ContextKey string

const (
	// ClaimsKey is the key used to store JWT claims in the context
	ClaimsKey ContextKey = "claims"
	// TenantIDKey is the key used to get tenant_id from claims
	TenantIDKey ContextKey = "tenant_id"
	// GinContextKey is the key used to store gin.Context in context.Context
	GinContextKey ContextKey = "GinContextKey"
)

var (
	ErrNoClaimsInContext   = errors.New("no claims found in context")
	ErrInvalidClaimsType   = errors.New("invalid claims type")
	ErrNoTenantIDInClaims  = errors.New("no tenant_id found in claims")
	ErrInvalidTenantIDType = errors.New("tenant_id must be a string")
	ErrNoGinContext        = errors.New("gin context not found")
)

func GetTenantIDFromContext(c context.Context) (string, error) {
	claims, exists := c.Value(string(ClaimsKey)).(jwt.MapClaims)
	if !exists {
		return "", ErrNoClaimsInContext
	}

	tenantID, exists := claims[string(TenantIDKey)]
	if !exists {
		return "", ErrNoTenantIDInClaims
	}

	tenantIDStr, ok := tenantID.(string)
	if !ok {
		return "", ErrInvalidTenantIDType
	}

	return tenantIDStr, nil
}
