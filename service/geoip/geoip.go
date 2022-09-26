package geoip

import (
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

var (
	db *geoip2.Reader
)

func Init(dbF string)  {
	var err error
	db, err = geoip2.Open(dbF)
	if err != nil {
		log.Fatal(err)
	}
}

func IsoCode(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	record, err := db.Country(ip)
	if err != nil {
		log.Printf("geoip err %s", err.Error())
		return ""
	}
	return record.Country.IsoCode
}
func Uninit()  {
	if db != nil {
		db.Close()
	}
}
