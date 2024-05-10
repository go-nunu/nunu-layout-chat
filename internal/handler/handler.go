package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"nunu-layout-chat/pkg/jwt"
	"nunu-layout-chat/pkg/log"
)

type Handler struct {
	logger *log.Logger
	component.Base
	app pitaya.Pitaya
}

func NewHandler(
	logger *log.Logger,
	app pitaya.Pitaya,
) *Handler {
	return &Handler{
		logger: logger,
		app:    app,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
