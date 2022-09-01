package controllers

import (
	"fmt"
	"line/models"
	"line/models/message"
	"line/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LineController struct {
	routerGroup *gin.RouterGroup
	service     *service.LineService
}

/*
建立物件
*/
func NewLineController(router *gin.Engine) *LineController {
	controller := LineController{
		router.Group("/line"),
		service.NewLineService(),
	}

	controller.receiveFromWebhook()
	controller.getMessageList()
	controller.sendMessage()
	controller.getUserList()

	return &controller
}

func (controller *LineController) receiveFromWebhook() {
	controller.routerGroup.POST("/webHook", func(c *gin.Context) {
		controller.service.ReceiveText(c.Writer, c.Request)
	})
}

// @Summary 取得訊息列表
// @Param index   path int true "第幾頁" default(1)
// @Param pageSize path int true "單頁顯示數量" default(10)
// @Param userId path string false "使用者代碼"
// @Router  /line/message/{index}/{pageSize} [GET]
// @Router  /line/message/{index}/{pageSize}/{userId} [GET]
func (controller *LineController) getMessageList() {
	controller.routerGroup.GET("/message/:index/:pageSize/*userId", func(c *gin.Context) {
		var pageInfo models.PageModel
		userId := c.Param("userId")
		userId = strings.Replace(userId, "/", "", 1)
		if err := c.ShouldBindUri(&pageInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		res := controller.service.GetMessageRecords(&pageInfo, userId)

		c.JSON(http.StatusOK, res)
	})
}

// @Summary 取得使用者列表
// @Param index   path int true "第幾頁" default(1)
// @Param pageSize path int true "單頁顯示數量" default(10)
// @Param name path string false "使用者名稱(關鍵字)" default()
// @Router  /line/user/{index}/{pageSize} [GET]
// @Router  /line/user/{index}/{pageSize}/{name} [GET]
func (controller *LineController) getUserList() {
	controller.routerGroup.GET("/user/:index/:pageSize/*name", func(c *gin.Context) {
		var pageInfo models.PageModel
		name := c.Param("name")
		name = strings.Replace(name, "/", "", 1)
		fmt.Println("name=" + name)
		if err := c.ShouldBindUri(&pageInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		res := controller.service.GetUserList(&pageInfo, name)

		c.JSON(http.StatusOK, res)
	})
}

// @Summary 發送訊息給使用者
// @accept application/json
// @Param request body message.SendMessageDto true "query params"
// @Router  /line/message [POST]
func (controller *LineController) sendMessage() {
	controller.routerGroup.POST("/message", func(c *gin.Context) {
		var json message.SendMessageDto
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		res := controller.service.SendMessage(json.UserId, json.Text)

		c.JSON(http.StatusOK, res)
	})
}
