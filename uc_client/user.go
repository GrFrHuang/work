package uc_client

import (
	"net/http"
	"log"
	"encoding/json"
)

type UserInfo struct {
	ID         int64
	UserName   string
	Nickname   string
	Email      string
	Password   string
	Extra      string
	Cellphone  string
	Registered string
	Open       string
}

func RegisterUser(name, pwd, email, phone, nickname string) (resname, resemail string, err error) {
	args := map[string]string{
		"UserName":name,
		"Email":email,
		"Password":pwd,
		"Cellphone":phone,
		"Nickname":nickname,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerUser, "POST", args); if err != nil {
		log.Printf("Response error:%v", err)
		return
	}
	var obj UserInfo
	if response.StatusCode == http.StatusOK {
		err = json.Unmarshal(response.Body, &obj); if err != nil {
			return
		}
		resname = obj.UserName
		resemail = obj.Email
		log.Printf("DecodeObj: %+v", obj)
		return
	} else {
		log.Printf("Response.StatusCode: %v", response.StatusCode)
		return "","", GetResErrInfo(response)
	}
}

func GetUserInfoByName(name, atoken string) (userInfo *UserInfo, err error) {
	args := map[string]string{
		"atoken":atoken,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerUser + "/" + name, "GET", args); if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		userInfo = new(UserInfo)
		err = json.Unmarshal(response.Body, userInfo); if err != nil {
			return
		}
		return
	} else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return nil, GetResErrInfo(response)
	}
}

func UpdateUserInfo(name, phone, extra, atoken string) (err error) {
	args := map[string]string{
		"UserName":name,
		"Extra":extra,
		"Cellphone":phone,
		"atoken":atoken,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerUser, "PUT", args); if err != nil {
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

func DeleteUser(name string) (err error) {
	args := map[string]string{
		"UserName":name,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerUser,"DELETE",args);if err != nil {
		log.Printf("response error:%v", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		return
	}else {
		log.Printf("response.StatusCode: %v", response.StatusCode)
		return GetResErrInfo(response)
	}
}

func ChangeOwnPwd(name, pwd, newpwd, atoken string) (err error) {
	args := map[string]string{
		"UserName":name,
		"Password":pwd,
		"NewPassword":newpwd,
		"atoken":atoken,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerChangePwd, "PUT", args); if err != nil {
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

func ChangeUserPwd(name, newpwd string) (err error) {
	args := map[string]string{
		"UserName":name,
		"NewPassword":newpwd,
	}
	response, err := RequestUcApi(ucenterUrl + apiVersion + routerChangeSomeonePwd, "PUT", args); if err != nil {
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


