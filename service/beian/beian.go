package beian

import (
	"database/sql"
	"github.com/1121170088/find-domain/search"
	"icp-search/dao"
	"icp-search/entity"
	init_ "icp-search/init"
	"icp-search/upstream"
	"icp-search/utils"
	"log"
	"strings"
	"sync"
)

var (
	checking = false
	mutex sync.Mutex
	limit = 100
	con chan bool
)

func CheckCode0(id int)  {
	if isChecking() {
		log.Printf("beian-checking has been running. ")
		return
	}
	setChecking()
	defer setUnChecking()
	if id != -1 {
		init_.Cfg.Code0Index = id
	}
back:
	icps, err := dao.SearchCode0FromId(init_.Cfg.Code0Index, limit)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if len(icps) == 0 {
		log.Printf("beian-checking finished")
		return
	}
	var icp *entity.Icp
	for _,  icp = range icps {
		con <- true
		go do(icp)
	}
	if icp != nil {
		init_.Cfg.Code0Index = icp.Id
	}
	goto back

}

func do(icp *entity.Icp)  {
	defer func() {
		<- con
	}()
	domain := icp.Domain
	nIcp, err := upstream.Mxnzp.Search(domain)
	if err != nil {
		return
	}
	icp.Code = 1
	icp.CacheTime = utils.CurrentDateTimeStr()
	icp.Name = nIcp.Name
	icp.PassTime = nIcp.PassTime
	icp.Type = nIcp.Type
	icp.IcpCode = nIcp.IcpCode
	icp.Unit = nIcp.Unit
	err = dao.Insert(icp)
	if err != nil {
		log.Printf(err.Error())
		return
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
	con = make(chan bool, init_.Cfg.Code0ConCurrent)
}

func Commit(domains []string)  {
	for _, domain := range domains {
		domain = strings.ToLower(domain)
		domain = search.Search(domain)
		if domain == "" {
			continue
		}
		con <- true
		go commit(domain)
	}
}

func commit(domain string)  {
	defer func() {
		<- con
	}()
	_, err := dao.Search(domain)
	if err != nil {
		if err == sql.ErrNoRows {
			var upStream upstream.Upstream
			upStream = upstream.Mxnzp
			icp, err := upStream.Search(domain)
			if err != nil {
				if err == upstream.Norecord {
					icp = &entity.Icp{
						Domain:    domain,
						Unit:      "",
						Type:      "",
						IcpCode:   "",
						Name:      "",
						PassTime:  "",
						CacheTime: utils.CurrentDateTimeStr(),
						Code:      0,
					}
					err := dao.Insert(icp)
					if err != nil {
						log.Printf("%v", err.Error())
					}
				}
				return
			} else {
				icp.Code = 1
				icp.CacheTime = utils.CurrentDateTimeStr()
				err := dao.Insert(icp)
				if err != nil {
					log.Printf("%v", err.Error())
				}
			}
		} else {
			log.Printf("%s", err.Error())
		}
	}
}
