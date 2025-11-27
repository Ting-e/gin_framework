package model

type GetListReq struct {
	Limit *int `form:"limit,omitempty"`
	Skip  *int `form:"skip,omitempty"`
}

type GetListResp struct {
	Code  int    `json:"code"`
	Mess  string `json:"mess"`
	Total int    `json:"total,omitempty"`
	Datas []User `json:"datas,omitempty"`
}
