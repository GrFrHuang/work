package tool

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego"
)

func TestSendInitPwdEmail(t *testing.T) {
	err := SendInitPwdEmail(map[string]string{
		"pwd":      "123456",
		"nickname": "nick",
		"name":     "nma",
	}, "zhangliang@kuaifazs.com")
	t.Error(err)
}

func Test(t *testing.T) {
	fmt.Println(beego.AppConfig.String("cps_url"))
	fmt.Println(beego.AppConfig.String("face_auth_host"))
}
