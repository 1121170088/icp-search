package upstream

import (
	"icp-search/entity"
	"io"
	"log"
	"net/http"
	"regexp"
)

var rgx = regexp.MustCompile(`<div class="box-title"><h2 class="btitle">备案信息</h2></div>\r\n\s*<div class=".*?">\r\n\s*<div class=".*?">主办单位性质</div>\r\n\s*<div class=".*?">(.*?)</div>\r\n\s*</div>\r\n\s*<div class=".*?">\r\n\s*<div class="line-name">主办单位名称</div>\r\n\s*<div class="line-value .*?">(.*?)</div>\r\n\s*</div>\r\n\s*<div class="line">\r\n\s*<div class="line-name .*?">主办单位备案号</div>\r\n\s*<div class="line-value">.*?</div>\r\n\s*</div>\r\n\s*<div class="line .*?">\r\n\s*<div class="line-name .*?">网站名称</div>\r\n\s*<div class="line-value">(.*?)</div>\r\n\s*</div>\r\n\s*<div class="line .*?">\r\n\s*<div class="line-name">网站备案号</div>\r\n\s*<div class="line-value .*?"><a href=".*?">(.*?)</a></div>\r\n\s*</div>\r\n\s*<div class="line .*?">\r\n\s*<div class="line-name">首页网址</div>\r\n\s*<div class="line-value"><a href=".*?" target="_blank" rel="nofollow">.*?</a></div>\r\n\s*</div>\r\n\s*<div class="line .*?">\r\n\s*<div class="line-name .*?">备案域名</div>\r\n\s*<div class="line-value">(.*?)</div>\r\n\s*</div>\r\n\s*<div class="line">\r\n\s*<div class="line-name">审核时间</div>\r\n\s*<div class="line-value .*?">(.*?)</div>\r\n\s*</div>`)
var rgx2 = regexp.MustCompile(`无备案信息`)
func Search(domain string) (*entity.Icp, error) {
	url := "https://www.beianhao.com/beian/" + domain
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return regexSearch(domain, bytes)
}

func regexSearch(domain string,bytes []byte) (*entity.Icp, error) {
	res := rgx.FindStringSubmatch(string(bytes))
	if len(res) == 7 {
		icp := &entity.Icp{
			Domain:   res[5],
			Unit:     res[2],
			Type:     res[1],
			IcpCode:  res[4],
			Name:     res[3],
			PassTime: res[6],
		}
		log.Printf("%v", icp)
		return icp, nil
	}
	res = rgx2.FindStringSubmatch(string(bytes))
	if len(res) == 1 {
		log.Printf("%s没查到", domain)
		return nil, Norecord
	}
	return nil, SearchErr

}