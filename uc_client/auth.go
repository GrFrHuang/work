package uc_client

import (
	"encoding/json"
	"log"
	"net/http"
)

type LoginResult struct {
	UserName             string
	RefreshToken         string
	AccessToken          string
	Session              string
	AccessTokenExpiresIn int
	SessionExpiresIn     int
}

// 登录
func Login(email, pwd string) (loginRet *LoginResult, err error) {
	args := map[string]string{
		"Email":    email,
		"Password": pwd,
	}
	response, err := RequestUcApi(ucenterUrl+apiVersion+routerCheckLogin, "POST", args)
	if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		loginRet = new(LoginResult)
		err = json.Unmarshal(response.Body, loginRet)
		if err != nil {
			return
		}
		return
	} else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return nil, GetResErrInfo(response)
	}
}

// 检查AccessToken
func CheckAccessToken(token, username string) (res bool, err error) {
	args := map[string]string{
		"atoken": token,
	}
	response, err := RequestUcApi(ucenterUrl+apiVersion+routerCheckAccessToken+"/"+username, "GET", args)
	if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		res = true
		return
	} else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return false, GetResErrInfo(response)
	}
}

// 通过RefreshToken 重置 AccessToken
func ResetAccessToken(rtoken, username string) (atoken string, err error) {
	args := map[string]string{
		"UserName":     username,
		"RefreshToken": rtoken,
	}
	response, err := RequestUcApi(ucenterUrl+apiVersion+routerResetAccessToken, "POST", args)
	if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		err = json.Unmarshal(response.Body, &atoken)
		if err != nil {
			return
		}
		return
	} else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return "", GetResErrInfo(response)
	}
}

func OffLine(name string) (err error) {
	args := map[string]string{}
	response, err := RequestUcApi(ucenterUrl+apiVersion+routerOffLine+"/"+name, "GET", args)
	if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		return
	} else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return GetResErrInfo(response)
	}
}
