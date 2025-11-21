package model

type GetListReq struct {
	Limit *int `form:"limit,omitempty"`
	Skip  *int `form:"skip,omitempty"`
}

type GetListResp struct {
	Code  int    `json:"code"`
	Mess  string `json:"mess"`
	Total int    `json:"total,omitempty"`
	Datas []Data `json:"datas,omitempty"`
}

type Data struct {
	Key         string `json:"key,omitempty"`
	Field       string `json:"field,omitempty"`
	Create_time string `json:"create_time,omitempty"`
}
