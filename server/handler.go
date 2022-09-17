package server

import (
	"database/sql"
	"encoding/json"
	"icp-search/dao"
	"icp-search/entity"
	init_ "icp-search/init"
	"icp-search/upstream"
	"icp-search/utils"
	"log"
	"net/http"
)
func Start()  {
	http.HandleFunc("/icp", func(writer http.ResponseWriter, request *http.Request) {
		domain := request.URL.Query().Get("domain")
		up := request.URL.Query().Get("up")
		header := writer.Header()
		header.Add("Content-Type", "application/json;charset=UTF-8")
		var err error
		resp := &entity.Resp{
			Code: 0,
			Msg:  "",
			Data: nil,
		}
		if domain == "" {
			resp.Msg = "域名不能为空"
			bytes, _ := json.Marshal(resp)
			writer.Write(bytes)
			return
		}
		icp, err := dao.Search(domain)
		if err != nil {
			if err == sql.ErrNoRows {
				var upStream upstream.Upstream
				switch up {
				case "1":upStream = upstream.Mxnzp
				case "2":upStream = upstream.BeiAnHao
				case "3":upStream = upstream.Miit
				default:
					upStream = upstream.Mxnzp
				}
				icp, err := upStream.Search(domain)
				if err != nil {
					if err == upstream.SearchErr {
						resp.Code = 2
						resp.Msg = "三方受限"
					} else if err == upstream.Norecord {
						resp.Code = 0
						resp.Msg = "可能没备案"
					} else {
						resp.Code = 2
						resp.Msg = "代码出错"
						log.Printf("%s", err.Error())
					}
				} else {
					resp.Data = icp
					resp.Code = 1
					resp.Msg = "ok"
					resp.Data.CacheTime = utils.CurrentDateStr()
					err := dao.Insert(icp)
					if err != nil {
						log.Printf("%v", err.Error())
					}
				}
			} else {
				resp.Msg = "未知数据库错误"
				log.Printf("%s", err.Error())
			}
		} else {
			resp.Code = 1
			resp.Msg = "ok"
			resp.Data = icp
		}
		bytes, _ := json.Marshal(resp)
		writer.Write(bytes)

	})
	http.ListenAndServe(init_.Cfg.Addr, nil)
}