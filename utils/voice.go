package utils

import (
	"net/http"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"github.com/astaxie/beego/toolbox"
	"errors"
	"fmt"
)

const (
	APIKEY    = "CPFONXAbnCw4n5pWYCkszWtR"
	SECRETKEY = "wqfYwY0e5z78UN6uwZUR2gMkkGSxuImk"
)

var token string

func VoiceRun() {
	getToken()
	// 每天4点获取新的token
	GetToken := toolbox.NewTask("GetToken", "0 0 4 * * * ", func() error {
		getToken()
		return nil
	})
	toolbox.AddTask("GetToken", GetToken)
	toolbox.StartTask()

}
func getToken() {
	var tokenUrl = "https://openapi.baidu.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + APIKEY + "&client_secret=" + SECRETKEY
	resp, err := http.Get(tokenUrl)
	checkErr(err)
	info, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	defer resp.Body.Close()
	json, err := simplejson.NewJson(info)
	checkErr(err)
	token, err = json.Get("access_token").String()
	fmt.Println("获取到百度的:"+token)
	checkErr(err)
}
func checkErr(err error) {
	if err != nil {
		print(err)
		return
	}
}

func VoiceByText(text string) (error, string) {
	if token == "" {
		return errors.New("请先调用VoiceRun方法"), ""
	}
	path := "http://tsn.baidu.com/text2audio?tex=" + text + "&lan=zh&cuid=123123&ctp=1&tok=" + token
	resp, err := http.Get(path)
	if err != nil {
		return err, ""
	}
	if resp.Header.Get("Content-Type") != "audio/mp3" {
		info, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err, ""
		}
		defer resp.Body.Close()
		return errors.New(string(info)), ""
	}
	return nil, path
}
