package utils

import (
	"github.com/astaxie/beego"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"errors"
)

var host = beego.AppConfig.String("face_auth_host")

func FaceVerify(uid string, base64 string) error {
	rep, err := http.PostForm("http://"+host+"/auth", url.Values{"uid": {uid}, "img": {base64}})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rep.Body.Close()
	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if string(data) != "ok" {
		return errors.New(string(data))
	}
	return nil
}
