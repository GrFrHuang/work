package controllers

import (
	"github.com/astaxie/beego"
)

// 修复数据
type RepairDataController struct {
	beego.Controller
}

////@router /repair [get]
//func (rc *RepairDataController) RepairData() {
//	err := models.ClearTable()
//	if err != nil {
//		beego.Warn("clear error : ", err)
//		return
//	}
//	accounts, _ := models.QueryChannelVerifyAccount()
//
//	for _, account := range accounts {
//		beego.Warning(account["GameStr"])
//		var gamea []models.GameAmount
//		err := json.Unmarshal([]byte(account["GameStr"].(string)), &gamea);
//		if err != nil {
//			beego.Warning(err.Error())
//		}
//		//将游戏id取到gid切片
//		for _, gameam := range gamea {
//			var m models.RepairData
//			m.GameId = gameam.GameId
//			m.Cp = account["Cp"].(string)
//			m.StartTime = account["StartTime"].(string)
//			m.EndTime = account["EndTime"].(string)
//			_, err := models.AddRepairData(&m);
//			if err != nil {
//				beego.Warning(err.Error())
//			}
//
//		}
//
//	}
//}
