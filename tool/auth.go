package tool

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/uc_client"
)

func GetRequestToken(ctx *context.Context) string {
	token := ctx.Input.Header("Authorization")
	if token == "" {
		token = ctx.Input.Query("_token")
	}
	return token
}

func VerifyToken(ctx *context.Context) (err error, uid int) {
	token := GetRequestToken(ctx)
	if token == "" || len(token) <= 5 {
		// debug 如果调试环境没有传token,就默认uid=1
		if beego.BConfig.RunMode == "dev" {
			uid = 1
			return
		}
		err = errors.New("forget token?")
		return
	}
	u, err := models.GetUserInfoByToken(token)
	if err != nil {
		err = errors.New("token invalid")
		return
	}
	res, err := uc_client.CheckAccessToken(token, u.Name)
	if err == nil && res {
		uid = u.Id
		return
	}
	return
}

func CreateSession(u bean.CreateSession) (loginRet *uc_client.LoginResult, err error) {
	loginRet, err = uc_client.Login(u.Email, u.Password)
	if err != nil {
		return
	}
	return
}


func OffLine(name string) (err error) {
	err = uc_client.OffLine(name)
	return
}
