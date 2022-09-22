package entity

type Icp struct {
	Id int `json:"-"`
	Domain string `json:"domain"`
	Unit string `json:"unit"`
	Type string `json:"type"`
	IcpCode string `json:"icpCode"`
	Name string `json:"name"`
	PassTime string `json:"passTime"`
	CacheTime string `json:"cacheTime"`
	Code int `json:"code"`
}
