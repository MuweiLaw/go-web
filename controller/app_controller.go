package controller

import (
	"github.com/gin-gonic/gin"
	"go-web/entity/dao"
	"go-web/entity/vo"
	"go-web/service"
	"net/http"
)

type appController struct {
	appService service.AppService
	//tokenService service.TokenService
}

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
		//app, err := ac.appService.Create(req)
		//if err != nil {
		//	return vo.InternalServerError.Data()
		//}
		var test = "ff"
		c.JSON(http.StatusOK, vo.Success(test))
	}
}
