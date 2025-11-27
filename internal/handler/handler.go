package handler

import (
	srv "project/internal/service"
)

var service srv.APIService

func init() {
	service = srv.GetService()
}
