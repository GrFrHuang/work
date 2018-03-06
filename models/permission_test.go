package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"testing"
	"time"
)

func TestPermission(t *testing.T) {
	link := "kftest:123456@(10.8.230.17:3307)/work_together"
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true

	uid := 74
	rsMap := GetAllRoleMap()
	user := User{Id:uid}
	orm.NewOrm().Read(&user)
	AddRoleInfo(&user, rsMap)
	log.Info("test", *user.Roles)
	x, e := CheckPermission(uid, "查", "contract_cp", nil)
	log.Info("TestPermission", "x:", x, "e:", e)

	time.Sleep(1000)
}



func TestMenu(t *testing.T) {
	uid := 62
	rsMap := GetAllRoleMap()
	user := User{Id:uid}
	orm.NewOrm().Read(&user)
	AddRoleInfo(&user, rsMap)
	log.Info("test", *user.Roles)
	x, e := GetCanVisitMenu(uid)
	log.Info("TestMenu", "menu:", x, "admin:", e)

	time.Sleep(1000)
}


func TestFilter(t *testing.T) {
	orm.Debug = true
	qs:=orm.NewOrm().QueryTable("user")
	u:=[]User{}
	qs.Filter("Id__in", 2,11).Filter("Id",1).All(&u)
	log.Info("x",u)
}
