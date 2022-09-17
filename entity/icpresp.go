package entity

type Resp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data *Icp `json:"data"`
}