package controller

import (
	"github.com/gin-gonic/gin"
	"go-web/entity/dao"
	"go-web/entity/vo"
	"go-web/service"
	"net/http"
)

var (
	as = *service.GetAppService()
)

func AppController(rg *gin.RouterGroup) {
	rg.POST("/create", create())
}

func create() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dao.CreateAppRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, vo.BadRequest.Data())
			return
		}
		app, err := as.Create(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Fail(err.Error(), nil))
			return
		}
		c.JSON(http.StatusOK, vo.Success(app))
	}
}
