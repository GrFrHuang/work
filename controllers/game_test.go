package controllers

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	"kuaifa.com/kuaifa/work-together/models"
)

// 修改CP商务负责人
// 游戏接入中接入人有，但对应发行商商务负责人没有，则将该人同步到发行商中
// 发行商中有，对应游戏接入也有，则以发行商中为准
func TestUpdateResponsiblePerson(t *testing.T) {
	o := orm.NewOrm()

	var companys []models.DistributionCompany
	o.Raw("SELECT * FROM distribution_company").QueryRows(&companys)
	for _, company := range companys {
		if company.YunduanResponsiblePerson == 0 {
			var game []models.Game
			o.Raw("SELECT DISTINCT(access_person) FROM game WHERE issue=? AND game_id>0 AND body_my=1",
				company.CompanyId).QueryRows(&game)
			if len(game) > 1 {
				o.Raw("SELECT DISTINCT(update_jruserid) FROM game WHERE issue=? AND game_id>0 AND body_my=1",
					company.CompanyId).QueryRows(&game)
				if len(game) > 1 {
					continue
				} else if len(game) == 1 {
					company.YunduanResponsiblePerson = game[0].UpdateJrUserID
					o.Update(&company)
				}
			} else if len(game) == 1 {
				company.YunduanResponsiblePerson = game[0].AccessPerson
				o.Update(&company)
			}
		}

		if company.YouliangResponsiblePerson == 0 {
			var game []models.Game
			o.Raw("SELECT DISTINCT(access_person) FROM game WHERE issue=? AND game_id>0 AND body_my=2",
				company.CompanyId).QueryRows(&game)
			if len(game) > 1 {
				o.Raw("SELECT DISTINCT(update_jruserid) FROM game WHERE issue=? AND game_id>0 AND body_my=2",
					company.CompanyId).QueryRows(&game)
				if len(game) > 1 {
					continue
				} else if len(game) == 1 {
					company.YouliangResponsiblePerson = game[0].UpdateJrUserID
					o.Update(&company)
				}
			} else if len(game) == 1 {
				company.YouliangResponsiblePerson = game[0].AccessPerson
				o.Update(&company)
			}
		}

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
