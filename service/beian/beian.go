package beian

import (
	"icp-search/dao"
	"icp-search/entity"
	init_ "icp-search/init"
	"icp-search/upstream"
	"icp-search/utils"
	"log"
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
