package test

import (
	"icp-search/service/geoip"
	"testing"
)

func TestIsoCode(t *testing.T)  {
	geoip.Init("../Country.mmdb")
	t.Log(geoip.IsoCode("192.168.1.1"))
	t.Log(geoip.IsoCode("49.51.78.122"))
	t.Log(geoip.IsoCode("128.14.151.195"))
	t.Log(geoip.IsoCode("107.150.102.158"))
	t.Log(geoip.IsoCode("108.167.146.59"))
	t.Log(geoip.IsoCode("69.28.62.189"))
	t.Log(geoip.IsoCode("68.178.224.4"))
	t.Log(geoip.IsoCode("34.102.136.180"))
	t.Log(geoip.IsoCode("47.91.170.222"))
	t.Log(geoip.IsoCode("160.20.59.185"))
	geoip.Uninit()
}
