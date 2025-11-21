package model

type AddDataReq struct {
	Field string `json:"field"`
}

type AddDataResp struct {
	Code int    `json:"code"`
	Mess string `json:"mess"`
}
