package model

type EditDataReq struct {
	ID    string `json:"id"`
	Field string `json:"field"`
}

type EditDataResp struct {
	Code int    `json:"code"`
	Mess string `json:"mess"`
}
