package checker

import (
	"github.com/astaxie/beego/orm"
	"testing"
	"time"
	"fmt"
)

func TestRemitAccount(t *testing.T) {
	err := RemitAccount()
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyAccount(t *testing.T) {
	err := VerifyAccount()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetNotVerifyAccount(t *testing.T) {
	r, err := GetAllNotVerifyChannel()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestContractTimeout(t *testing.T) {
	err := ContractTimeout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestContractNotSign(t *testing.T) {
	err := ContractNotSign()
	if err != nil {
		t.Fatal(err)
	}
}

func init() {
	link := "kftest:123456@(10.8.230.17:3308)/work_together"
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}

func TestTime(t *testing.T) {
	str := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	fmt.Println(str)
}
