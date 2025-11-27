package handler

import (
	"encoding/json"
	"project/examples/simple_crud/model"
	srv "project/examples/simple_crud/service"
	"project/pkg/logger"
	"project/pkg/response"

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
		response.BadRequest(c)
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("GetDatas （入参）:", string(req))

	//调用service层
	res := service.GetList(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("GetDatas （出参）:", string(resp))
	response.Success(c, res)
}

func GetData(c *gin.Context) {

	ID := c.Param("id")
	if ID == "" {
		logger.Sugar.Error("[controller] GetData: error: Not Found Param")
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("GetDatas （入参）:", ID)

	//调用service层
	res := service.GetData(ID)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("GetDatas （出参）:", string(resp))
	response.Success(c, res)
}

func AddData(c *gin.Context) {
	var body model.AddDataReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Sugar.Error("[controller] AddData: parse data error: ", err)
		response.BadRequest(c)
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("AddData （入参）:", string(req))

	//调用service层
	res := service.AddData(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("AddData （出参）:", string(resp))
	response.Success(c, res)
}

func DelData(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		logger.Sugar.Error("[controller] DelData: error: not found param")
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("DelData （入参）:", ID)

	//调用service层
	res := service.DelData(ID)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("DelData （出参）:", string(resp))
	response.Success(c, res)
}

func EditData(c *gin.Context) {
	var body model.EditDataReq

	body.ID = c.Param("id")
	if body.ID == "" {
		logger.Sugar.Error("[controller] EditData: error: not found param")
		response.BadRequest(c)
		return
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Sugar.Error("[controller] EditData: parse data error: ", err)
		response.BadRequest(c)
		return
	}

	req, _ := json.Marshal(body)
	logger.Sugar.Info("EditData （入参）:", string(req))

	//调用service层
	res := service.EditData(body)
	resp, _ := json.Marshal(res)
	logger.Sugar.Info("EditData （出参）:", string(resp))
	response.Success(c, res)
}
