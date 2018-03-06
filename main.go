package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "kuaifa.com/kuaifa/work-together/routers"
	"kuaifa.com/kuaifa/work-together/task"
	"kuaifa.com/kuaifa/work-together/utils"
	"kuaifa.com/kuaifa/work-together/controllers"
)

func main() {
	//models.Syncdb()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//Add kpi task listener
	go func() {
		controllers.ExecuteListenTask()
	}()

	task.Run()
	utils.VoiceRun()
	beego.BConfig.Listen.ServerTimeOut = 10
	beego.Run()
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
		beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = beego.BConfig.RunMode == "dev"
	//logs.SetLogger(logs.AdapterFile, `{"filename":"info.log","level":6,"maxlines":0,"maxsize":0,"daily":true,"maxdays":1000}`)
}
