package controller

import (
	"github.com/gin-gonic/gin"
)

type controller interface {
	New(e *gin.Engine)
}
