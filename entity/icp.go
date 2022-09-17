package entity

type Icp struct {
	Domain string `json:"domain"`
	Unit string `json:"unit"`
	Type string `json:"type"`
	IcpCode string `json:"icpCode"`
	Name string `json:"name"`
	PassTime string `json:"passTime"`
	CacheTime string `json:"cacheTime"`
}
