package old_sys

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/http_util"
	"strconv"
	"time"
	"log"
)

const (
	SECRET_KEY = "5f5a740a75ed31c67c5d16eded29d30d"
	HOST       = "http://www.kuaifazs.com"
)

func sign(parms string) string {
	m := md5.New()
	m.Write([]byte(parms))
	tmp_md5 := hex.EncodeToString(m.Sum(nil))
	mm := md5.New()
	mm.Write([]byte(tmp_md5 + SECRET_KEY))
	return hex.EncodeToString(mm.Sum(nil))
}

// params: 除了timestamp 和 sign 之外的参数
func httpRequest(api string, params map[string]string) (response string, err error) {
	url := HOST + api
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params["timestamp"] = timestamp
	p := util.ParseOrderKV(params)
	p.Sort()
	sign := sign(p.EncodeStringWithoutEscape())
	p.Add("sign", sign)
	_, response, err = http_util.Get(url, p, nil)

	log.Printf("url: %s?%s", url, p.EncodeString())
	log.Printf("rsp: %s", response)
	return
}

// params: 除了timestamp 和 sign 之外的参数
func request(api string, params map[string]string, rsp interface{}) (err error) {
	response, err := httpRequest(api, params)
	if err != nil {
		return
	}
	bj, err := bjson.New([]byte(response))
	if err != nil {
		return
	}
	if bj.Pos("status").Int() != 0 {
		err = errors.New(bj.Pos("error_message").String())
		return
	}
	if rsp != nil {
		err = bj.Pos("data").Object(rsp)
		if err != nil {
			return
		}
	}

	return
}

// 阶梯规则 ladder
