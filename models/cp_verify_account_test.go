package models

import (
	"testing"
	"strings"
	"time"
)

//func TestGetGameUpToDateTime(t *testing.T) {
//	uptime, err := GetGameUpToDateTime(0); if err != nil {
//		t.Fatalf("error: %v", err)
//	}
//	t.Logf("uptime: %v", uptime)
//}

func TestName(t *testing.T) {
	//GetNoCpVerifyAccount()
}

func TestSumAmonut(t *testing.T) {
	//res := SumAmonut()
	//t.Logf("res: %v",res)

	create := strings.Repeat(",?",3)[1:]
	//create = create[1:]
	t.Logf("res: %v",create)
	t.Logf("time: %v",time.Now().Unix())
}

func TestGetVerifyDateByGameid(t *testing.T) {
	//GetVerifyDateByGameid(2)
}

func TestDate(t *testing.T) {
	s:= "2016-12"
	tm, err := time.Parse("2006-01",s); if err != nil {
		t.Fatalf("error: %v",err)
	}
	t.Logf("res: %v",tm.AddDate(0,1,-1))
}


//func init() {
//	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
//		"123456", "10.8.230.17",
//		"3307", "work_together")
//	fmt.Printf("link:%v", link)
//	orm.RegisterDataBase("default", "mysql", link)
//
//	orm.Debug = true
//}
