package checker

import (
	"github.com/astaxie/beego/orm"
	"testing"
	"time"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models"
	"reflect"
)

//func TestContractNotSign(t *testing.T) {
//	err := ContractNotSign()
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func init() {
	//link := "kftest:123456@(10.8.230.17:3308)/work_together"
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
		"kuaifazs", "10.8.230.17",
		"3308", "work_together_online")
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
	//	"123456", "10.8.230.17",
	//	"3308", "work_together")
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}

func TestGetEmailsByUserIds(t *testing.T) {
	rules, _ := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_PAY)
	tmp := GetEmailsByUserIds(rules[0].UserIds)
	fmt.Println(tmp)
}

func TestTime(t *testing.T)  {
	str := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	fmt.Println(str)
}

func TestChannelContractExpireWarning(t *testing.T) {
	err := ChannelContractExpireWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestNewGameAccessWarning(t *testing.T) {
	err := SendNewGameAccessWarning(267)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCpContractExpireWarning(t *testing.T) {
	err := CpContractExpireWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestChannelContractSignWarning(t *testing.T) {
	err := ChannelContractSignWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestChannelOrderVerifyWarning(t *testing.T) {
	err := ChannelOrderVerifyWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCpOrderVerifyWarning(t *testing.T) {
	err := CpOrderVerifyWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestChannelOrderThresholdWarning(t *testing.T) {
	err := ChannelOrderThresholdWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestChannelOrderPayWarning(t *testing.T) {
	err := ChannelOrderPayWarning()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestSendNewGameAccessWarning(t *testing.T) {
	err := SendNewGameAccessWarning(207)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestReflect(t *testing.T){
	//old := nil
	new := models.Contract{}
	//oldValue := reflect.Indirect(reflect.ValueOf(nil))
	newValue := reflect.Indirect(reflect.ValueOf(new))

	fmt.Printf("reflect.ValueOf(new):%v---\n", reflect.ValueOf(new))
	//fmt.Printf("oldTypes:%v---\n", oldValue)
	fmt.Printf("newTypes:%v---\n", newValue)

	//oldTypes := oldValue.Type()
	newTypes := newValue.Type()

	//fmt.Printf("oldTypes:%v---\n", oldTypes)
	fmt.Printf("newTypes:%v---\n", newTypes)


}