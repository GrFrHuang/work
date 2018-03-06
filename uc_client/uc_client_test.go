package uc_client

import "testing"

var (
	username = "test001"
	pwd      = "654321"
	email    = "navery@foxmail.com"
)

func TestRegisterUser(t *testing.T) {
	resname, _, err := RegisterUser(username, pwd, email, "", "")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("Reg username: %v", resname)
}

func TestLogin(t *testing.T) {
	loginRet, err := Login(email, pwd)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("LoginRet %+v", loginRet)
}

var token = "3fdd5b604516ef799de1f42013add649"

func TestCheckAccessToken(t *testing.T) {
	res, err := CheckAccessToken(token, username)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if res {
		t.Log("success")
	}
}

func TestResetAccessToken(t *testing.T) {
	rtoken := "5805ce2407f8df3861a4c5ec751c203f"
	atoken, err := ResetAccessToken(rtoken, username)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("ResAtoken: %v", atoken)
}

func TestUpdateUserInfo(t *testing.T) {
	err := UpdateUserInfo(username, "15528297560", "hahaha", token)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Log("success!")
}

func TestChangeUserPwd(t *testing.T) {
	err := ChangeUserPwd(username, "654321")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Log("success!")
}

func TestGetUserInfoByName(t *testing.T) {
	u, err := GetUserInfoByName(username, token)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Log("success!")
	t.Logf("UserInfo %+v", u)
}

func TestOffLine(t *testing.T) {
	err := OffLine(username)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Log("success!")
}
