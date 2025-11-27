package model

type EditDataReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EditDataResp struct {
	Code int    `json:"code"`
	Mess string `json:"mess"`
}
