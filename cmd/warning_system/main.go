package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	//_ "github.com/go-sql-driver/mysql"
	"kuaifa.com/kuaifa/work-together/cmd/warning_system/checker"
	"kuaifa.com/kuaifa/work-together/cmd/contract_state/change"
)

func main() {
	walk()
}

func walk() {
	// 渠道合同过期
	err := checker.ChannelContractExpireWarning()
	if err != nil {
		log.Error("WARNING", "ChannelContractSignWarning", err)
	}

	// CP合同过期
	err = checker.CpContractExpireWarning()
	if err != nil {
		log.Error("WARNING", "CpContractExpireWarning", err)
	}

	// 渠道合同未签订
	err = checker.ChannelContractSignWarning()
	if err != nil {
		log.Error("WARNING", "ChannelContractSignWarning", err)
	}

	// 渠道未对账
	err = checker.ChannelOrderVerifyWarning()
	if err != nil {
		log.Error("WARNING", "ChannelOrderVerifyWarning", err)
	}

	// CP未对账
	err = checker.CpOrderVerifyWarning()
	if err != nil {
		log.Error("WARNING", "CpOrderVerifyWarning", err)
	}

	// 游戏流水阈值
	err = checker.ChannelOrderThresholdWarning()
	if err != nil {
		log.Error("WARNING", "ChannelOrderThresholdWarning", err)
	}

	// 渠道未回款
	err = checker.ChannelOrderPayWarning()
	if err != nil {
		log.Error("WARNING", "ChannelOrderPayWarning", err)
	}

	//判断合同状态是否到期，并修改合同状态
	err = change.ChangeContractState()
	if err != nil{
		log.Error("WARNING", "ChangeContractState", err)
	}

	// 由于此处的预警需要单独编译并使用crontab定时调度，为了简便，新加的预警如：
	// 游戏停运公告、停运游戏CP合同待处理、停运游戏渠道合同待处理、停运游戏CP合同未下架、停运游戏渠道合同未下架等
	// 放在了task定时任务中，使用beego的定时toolbox
}


func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
		beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}
