package handler

import (
	"encoding/json"
	"net/http"
	"project/api"
	"project/internal/model"
	srv "project/internal/service"
	"project/pkg/logger"

	"github.com/gin-gonic/gin"
)

var service api.APIService

func init() {
	service = srv.GetService()
}

func GetData(c *gin.Context) {
	var body model.GetDatasReq
	err := c.ShouldBind(&body)
	if err != nil {
		logger.Sugar.Error("[controller] GetDatas: parse data error: ", err)
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数解析出错"})
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("GetDatas （入参）:", string(req))

	//调用service层
	res := service.GetDatas(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("GetDatas （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}
