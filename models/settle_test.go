package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"testing"
)


//func TestGetNotSettleAccount(t *testing.T) {
//	x, e := GetNotSettleAccount([]interface{}{},1485792001,0)
//	if e != nil {
//		//t.Error(e)
//	}
//	log.Info("test", fmt.Sprintf("%+v", x), e)
//}
//
//func TestDoPushSettleToVerify(t *testing.T) {
//	err:=DoPushSettleToVerify(1)
//	log.Info("test",err)
//}
//
//func TestDoPushRemitToVerify(t *testing.T) {
//	err:=DoPushRemitToVerify(0)
//	log.Info("test",err)
//}

func TestOrm(t *testing.T) {
	o:=orm.NewOrm()
	p:=[]orm.Params{}
	_,err:=o.QueryTable("order").Filter("name", "slene").Values(&p,"MAX(time)")
	if err != nil {
		log.Info("test err: ",err)
		return
	}
	log.Info("test p: ",p)
}

