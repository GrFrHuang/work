package task

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	"kuaifa.com/kuaifa/work-together/task/warning"
	"time"
)

// 测试利润计算
func TestComputeProfit(t *testing.T) {
	//start := time.Now()
	start, _ := time.Parse("2006-01-02", "2018-02-22")
	end, _ := time.Parse("2006-01-02", "2018-02-22")
	for i := start; i.Before(end.AddDate(0, 0, 1)); i = i.AddDate(0, 0, 1) {
		tim := i.Format("2006-01-02")
		ComputeProfit(tim)
	}
}

// 测试将通讯录渠道负责人同步到渠道接入中
func TestSyncBusiness(t *testing.T) {
	//FirstSyncChannelBusiness()
	//SyncChannelBusiness()

	//FirstSyncGameAccessPerson()
	SyncGameAccessPerson()
}

// 测试更新主合同状态为已到期
func TestChangeMainContractState(t *testing.T) {
	ChangeMainContractState()
}

// 测试根据流水生成合同
func TestAddContractByOrder(t *testing.T) {
	AddContractByOrder()
}

func TestCompute(t *testing.T) {
	//var incrTime int64 = 1510225200
	//incrDays := (incrTime - time.Now().Unix()) / 3600 / 24
	//fmt.Printf("day:%v\n", incrDays)

	intervals := []int{0, 1, 2}
	dates := warning.IntervalsToDates(false, intervals...)
	fmt.Printf("day:%v\n", dates)

	//day := time.Now().AddDate(0, 0, 0).Format("2006-01-02")
	//fmt.Printf("day:%v\n", day)
	//
	//unix := time.Unix(incrTime, incrTime)
	//da := warning.TimeSubDays(unix, time.Now())
	//fmt.Printf("unix:%v,day:%v\n", unix, da)

}

// 测试所有游戏停运相关的
func TestAllGameoutage(t *testing.T) {
	warning.DownContract()
	warning.GameOutageChannelContractOpe()
	warning.GameOutageCpContractOpe()
	warning.GameOutageChannelContractDown()
	warning.GameOutageCpContractDown()
}

// 测试游戏停运
func TestGameoutage(t *testing.T) {
	err := warning.GameOutageWarning()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

// 测试游戏合同下架
func TestDownContract(t *testing.T) {
	err := warning.DownContract()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

// 测试游戏停运渠道合同待处理
func TestGameOutageChannelContractOpe(t *testing.T) {
	err := warning.GameOutageChannelContractOpe()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

// 测试游戏停运cp合同待处理
func TestGameOutageCpContractOpe(t *testing.T) {
	err := warning.GameOutageCpContractOpe()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

// 测试停运游戏渠道合同未下架
func TestGameOutageChannelContractDown(t *testing.T) {
	err := warning.GameOutageChannelContractDown()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

// 停运游戏CP合同未下架
func TestGameOutageCpContractDown(t *testing.T) {
	err := warning.GameOutageCpContractDown()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

func TestMainContract(t *testing.T) {
	// 渠道主合同过期预警
	err := warning.ChannelMainContractExpireWarning()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}

	// CP主合同过期预警
	err = warning.CpMainContractExpireWarning()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
}

func init() {
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
	//	"kuaifazs", "10.8.230.17",
	//	"3308", "work_together_online")
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
		"123456", "10.8.230.17",
		"3308", "work_together")
	fmt.Printf("link:%v", link)
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
