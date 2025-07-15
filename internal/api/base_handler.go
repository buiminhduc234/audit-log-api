package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

func (h *BaseHandler) RequestCtx(ginCtx *gin.Context) context.Context {
	ctx := ginCtx.Request.Context()
	for k, v := range ginCtx.Keys {
		ctx = context.WithValue(ctx, k, v) // nolint: staticcheck
	}
	return ctx
}
