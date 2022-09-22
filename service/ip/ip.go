package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/babolivier/go-doh-client"
	"icp-search/dao"
	init_ "icp-search/init"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	resolver doh.Resolver
	checking = false
	mutex sync.Mutex
	limit = 100
)

func CheckIp()  {
	if isChecking() {
		log.Printf("ip-checking has been running. ")
		return
	}
	setChecking()
	defer setUnChecking()
back:
	icps, err := dao.SearchById(init_.Cfg.IpIndex, limit)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if len(icps) == 0 {
		log.Printf("ip-checking finished")
		return
	}
	for _, icp := range icps {
		domain := icp.Domain
		var ip string
		ip, err = lookup(domain)
		if err != nil {
			ip, err = lookup("www." + domain)
		}
		if err == nil {
			icp.Ip = ip
			isoCode := searchIsoCode(ip)
			icp.IsoCode = isoCode
			err = dao.Insert(icp)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			init_.Cfg.IpIndex = icp.Id
		} else {
			msg := fmt.Sprintf("looking up ip err %s %s", domain, err.Error())
			log.Printf(msg)
		}
		err = nil
	}
	if err == nil {
		goto back
	}

}

func isChecking() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return checking
}

func setChecking()  {
	mutex.Lock()
	defer mutex.Unlock()
	checking = true
}

func setUnChecking()  {
	mutex.Lock()
	defer mutex.Unlock()
	checking = false
}
func init()  {
	resolver = doh.Resolver{
		Host:  init_.Cfg.Doh, // Change this with your favourite DoH-compliant resolver.
		Class: doh.IN,
	}
}

func lookup(domain string) (ip string, err error)  {
	a, _, err := resolver.LookupA(domain)
	if err != nil {
		return "", err
	}
	if len(a) == 0 {
		return "", errors.New("invalid lookup")
	}
	ip = a[0].IP4
	return ip, nil
}

func searchIsoCode(ip string) string  {
	response, err := http.Get(init_.Cfg.IpInfoServer + ip)
	if err != nil {
		msg := fmt.Sprintf("ip info server err -> ip: %s, err: %s", ip, err.Error())
		log.Printf(msg)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	var info = &struct {
		Code string
	}{}
	err = json.Unmarshal(body, info)
	if err != nil {
		log.Printf(err.Error())
	}
	return info.Code
}