package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	_"github.com/go-sql-driver/mysql"
	"kuaifa.com/kuaifa/work-together/cmd/alarm_system/checker"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

const interval = time.Hour * 72

// 报警系统
// 每 ~8h 去系统中找应该报警东西
func main() {
	for {
		walk()
		time.Sleep(interval)
	}
}

func walk() {
	err := clean()
	if err != nil {
		log.Error("Alarm", "clean", err)
	}

	err = checker.RemitAccount()
	if err != nil {
		log.Error("Alarm", "RemitAccount", err)
	}
	err = checker.VerifyAccount()
	if err != nil {
		log.Error("Alarm", "VerifyAccount", err)
	}
	err = checker.ContractTimeout()
	if err != nil {
		log.Error("Alarm", "ContractTimeout", err)
	}
	err = checker.ContractNotSign()
	if err != nil {
		log.Error("Alarm", "ContractNotSign", err)
	}
}

// 清除以前的警告
func clean() (err error) {
	hide := map[string]interface{}{
		"IsHide":1,
	}

	o := orm.NewOrm()
	_, err = o.QueryTable(new(models.AlarmLog)).Update(hide)
	return
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
		beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}
