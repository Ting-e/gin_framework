package model

type AddDataReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type AddDataResp struct {
	Code int    `json:"code"`
	Mess string `json:"mess"`
}
