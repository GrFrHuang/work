package change

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func init() {
	link := "kftest:123456@(10.8.230.17:3308)/work_together"
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}

func TestChangeContractState(t *testing.T){
	err := ChangeContractState()
	if err != nil{
		fmt.Println(err.Error())
	}
}
