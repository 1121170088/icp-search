package upstream

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"icp-search/entity"
	init_ "icp-search/init"
	"io"
	"log"
	"net/http"
	"time"
)
var SearchErr = errors.New("search error")
var Norecord = errors.New("no recrod")

var errored = false
var date string
type MxnzpUpstream struct{}
var Mxnzp = &MxnzpUpstream{}
func (u *MxnzpUpstream) Search(domain string) (*entity.Icp, error) {

	if errored {
		ndata := time.Now().Format(time.Now().Format("20060102"))
		if ndata != date {
			errored = false
		} else {
			return nil, SearchErr
		}
	}

	b64d:= base64.StdEncoding.EncodeToString([]byte(domain))
	url := fmt.Sprintf(init_.Cfg.Upstream, b64d)
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:103.0) Gecko/20100101 Firefox/103.0")
	reqest.Header.Add("Pragma", "no-cache")
	reqest.Header.Add("Cache-Control", "Cache-Control")
	reqest.Header.Add("Accept", "application/json")

	resp, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf(string(bytes))
	res := &entity.Resp{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}
	if res.Code == 1 {
		return res.Data, nil
	}
	if res.Code == 0 {
		return nil, Norecord
	}
	errored = true
	date = time.Now().Format(time.Now().Format("20060102"))
	return nil, SearchErr
}
