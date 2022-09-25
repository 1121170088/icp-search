package server

import (
	"database/sql"
	"encoding/json"
	"github.com/1121170088/find-domain/search"
	"icp-search/dao"
	"icp-search/entity"
	init_ "icp-search/init"
	"icp-search/service/beian"
	"icp-search/service/ip"
	"icp-search/upstream"
	"icp-search/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		domain = strings.ToLower(domain)
		domain = search.Search(domain)
		if domain == "" {
			resp.Msg = "可能不是域名"
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
						resp.Data = icp
						resp.Code = icp.Code
						resp.Msg = "可能没备案"
						err := dao.Insert(icp)
						if err != nil {
							log.Printf("%v", err.Error())
						}
					} else {
						resp.Code = 2
						resp.Msg = "代码出错"
						log.Printf("%s", err.Error())
					}
				} else {
					icp.Code = 1
					icp.CacheTime = utils.CurrentDateTimeStr()
					resp.Data = icp
					resp.Code = icp.Code
					resp.Msg = "ok"
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
			resp.Code = icp.Code
			resp.Msg = "ok"
			if icp.Code == 0 {
				resp.Msg = "可能没备案"
			}
			resp.Data = icp
		}
		bytes, _ := json.Marshal(resp)
		writer.Write(bytes)

	})
	http.HandleFunc("/check0", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		var id = -1
		idstr := request.URL.Query().Get("id")
		if idstr != "" {
			id, err = strconv.Atoi(idstr)
			if err != nil {
				writer.Write([]byte("fail"))
				return
			}
		}
		go beian.CheckCode0(id)
		writer.Write([]byte("ok"))
	})
	http.HandleFunc("/checkIp", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		var id = -1
		idstr := request.URL.Query().Get("id")
		if idstr != "" {
			id, err = strconv.Atoi(idstr)
			if err != nil {
				writer.Write([]byte("fail"))
				return
			}
		}
		go ip.CheckIp(id)
		writer.Write([]byte("ok"))
	})
	http.ListenAndServe(init_.Cfg.Addr, nil)
}