package model

type GetDatasReq struct {
	Key   string `json:"key,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Skip  int    `json:"skip,omitempty"`
}
type GetDatasResp struct {
	Code  int     `json:"code"`
	Mess  string  `json:"mess"`
	Total int     `json:"total,omitempty"`
	Datas []*Data `json:"datas,omitempty"`
}

type Data struct {
	Key         string `json:"key,omitempty"`
	Field       string `json:"field,omitempty"`
	Create_time string `json:"create_time,omitempty"`
}
