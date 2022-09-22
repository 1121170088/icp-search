package service

import (
	"icp-search/dao"
	"icp-search/entity"
	"icp-search/upstream"
	"icp-search/utils"
	"log"
	"sync"
)

var (
	checking = false
	mutex sync.Mutex

	id = 0
	code = 0
	limit = 100
)

func CheckCode0()  {
	if isChecking() {
		log.Printf("checking has been running. ")
		return
	}
	setChecking()
	defer setUnChecking()
back:
	icps, err := dao.SearchByCodeId(id, code, limit)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if len(icps) == 0 {
		log.Printf("checking finished")
		return
	}
	for _, icp := range icps {
		var nIcp *entity.Icp
		nIcp, err = upstream.Mxnzp.Search(icp.Domain)
		if err != nil {
			if err == upstream.Norecord {
				icp.CacheTime = utils.CurrentDateTimeStr()
				nIcp = icp
				err = nil
			} else {
				log.Printf(err.Error())
				break
			}
		}
		err = dao.Insert(nIcp)
		if err != nil {
			log.Printf(err.Error())
			break
		}
		id = nIcp.Id
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