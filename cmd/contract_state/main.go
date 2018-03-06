package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	//_ "github.com/go-sql-driver/mysql"
	"kuaifa.com/kuaifa/work-together/cmd/contract_state/change"
)

func main() {
	walk()
}

func walk() {
	//判断合同状态是否到期，并修改合同状态
	err := change.ChangeContractState()
	if err != nil{
		log.Error("WARNING", "ChangeContractState", err)
	}
}


func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
		beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}
