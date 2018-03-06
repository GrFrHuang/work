package change

import (
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"github.com/bysir-zl/bygo/log"
)

func ChangeContractState()(err error){

	o := orm.NewOrm()

	var cons []models.Contract

	timeNow := time.Now().Unix()
	qs := o.QueryTable(new(models.Contract)).Filter("state__exact", 150).Filter("end_time__lte", timeNow)

	_, err = qs.All(&cons)
	if(err != nil){
		return
	}


	for _, con := range cons{
		con.State = 151
		if _, err = o.Update(&con); err != nil{
			return
		}else {
			log.Info("INFO", "Update ContractState Success,id = " , con.Id)
		}
		//fmt.Printf("state:%v----timeNow:%v--------endTime:%v--------id:%v\n", con.State, timeNow, con.EndTime, con.Id)
	}

	return
}