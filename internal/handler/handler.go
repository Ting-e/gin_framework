package handler

import (
	"encoding/json"
	"net/http"
	"project/internal/model"
	srv "project/internal/service"
	"project/pkg/logger"

	"github.com/gin-gonic/gin"
)

var service srv.APIService

func init() {
	service = srv.GetService()
}

func GetList(c *gin.Context) {
	var body model.GetListReq
	err := c.ShouldBind(&body)
	if err != nil {
		logger.Sugar.Error("[controller] GetDatas: parse data error: ", err)
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数解析出错"})
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("GetDatas （入参）:", string(req))

	//调用service层
	res := service.GetList(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("GetDatas （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}

func GetData(c *gin.Context) {

	ID := c.Param("id")
	if ID == "" {
		logger.Sugar.Error("[controller] GetData: error: not found param")
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数缺失"})
		return
	}

	logger.Sugar.Info("GetDatas （入参）:", ID)

	//调用service层
	res := service.GetData(ID)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("GetDatas （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}

func AddData(c *gin.Context) {
	var body model.AddDataReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Sugar.Error("[controller] AddData: parse data error: ", err)
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数解析出错"})
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("AddData （入参）:", string(req))

	//调用service层
	res := service.AddData(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("AddData （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}

func DelData(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		logger.Sugar.Error("[controller] DelData: error: not found param")
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数缺失"})
		return
	}

	logger.Sugar.Info("DelData （入参）:", ID)

	//调用service层
	res := service.DelData(ID)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("DelData （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}

func EditData(c *gin.Context) {
	var body model.EditDataReq

	body.ID = c.Param("id")
	if body.ID == "" {
		logger.Sugar.Error("[controller] EditData: error: not found param")
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数缺失"})
		return
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Sugar.Error("[controller] EditData: parse data error: ", err)
		c.JSON(http.StatusBadRequest, model.Response{Code: http.StatusBadRequest, Message: "参数解析出错"})
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("EditData （入参）:", string(req))

	//调用service层
	res := service.EditData(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("EditData （出参）:", string(resp))
	c.JSON(http.StatusOK, res)
}
