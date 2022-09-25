package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/babolivier/go-doh-client"
	"icp-search/dao"
	"icp-search/entity"
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
	con chan bool
)

func CheckIp(id int)  {
	if isChecking() {
		log.Printf("ip-checking has been running. ")
		return
	}
	setChecking()
	defer setUnChecking()
	if id != -1 {
		init_.Cfg.IpIndex = id
	}
back:
	icps, err := dao.SearchFromId(init_.Cfg.IpIndex, limit)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if len(icps) == 0 {
		log.Printf("ip-checking finished")
		return
	}
	var icp *entity.Icp
	for _,  icp = range icps {
		con <- true
		go do(icp)
	}
	if icp != nil {
		init_.Cfg.IpIndex = icp.Id
	}
	goto back

}

func do(icp *entity.Icp)  {
	defer func() {
		<- con
	}()
	domain := icp.Domain
	ip, err := lookup(domain)
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
	} else {
		msg := fmt.Sprintf("looking up ip err %s %s", domain, err.Error())
		log.Printf(msg)
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
	con = make(chan bool, init_.Cfg.IpCheckConCurrent)
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