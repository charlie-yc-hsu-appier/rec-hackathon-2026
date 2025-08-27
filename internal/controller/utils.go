package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleBadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "detail": err.Error()})
}
