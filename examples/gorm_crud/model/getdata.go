package model

type GetDataReq struct {
	ID string
}

type GetDataResp struct {
	Code  int    `json:"code"`
	Mess  string `json:"mess"`
	Datas User   `json:"datas"`
}
