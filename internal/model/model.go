package model

type Resp struct {
	Code int    `json:"code"`
	Mess string `json:"mess,omitempty"`
}

type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"mess,omitempty"`
}
