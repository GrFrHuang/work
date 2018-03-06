package task

import (
	"time"
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"encoding/json"
	"github.com/astaxie/beego"
)

// 根据流水自动生成渠道合同
// 公司内部的快发、半袋米和好玩这几个渠道已经运营且有流水的游戏自动同步到WT的渠道合同里，以便查询平台数据的连续性。
func AddContractByOrder() {
	date := time.Now().Format("2006-01-02")
	fmt.Printf("date:%v\n", date)

	o := orm.NewOrm()
	var orders []models.Order
	o.QueryTable(new(models.Order)).Filter("date__exact", date).Filter("cp__in", "kuaifa", "bandaimi", "family").
		All(&orders, "game_id", "date", "cp")

	if len(orders) != 0 {
		// 以上三个渠道如果有流水，则循环判断是否有渠道合同，如果没有则生成
		for _, order := range orders {
			flag := o.QueryTable(new(models.Contract)).Filter("company_type__exact", 1).Filter("game_id__exact", order.GameID).
				Filter("channel_code__exact", order.Cp).Exist()
			if flag == false {
				// 没有渠道合同需要生成

				// 渠道合同默认分成为5 5 0
				ladder := []old_sys.Ladder4Post{{
					SlottingFee: 0,
					Ratio:       0.5,
				}}
				byteLadder, _ := json.Marshal(ladder)

				var con models.Contract
				con.State = 149
				con.GameId = order.GameID
				con.ChannelCode = order.Cp
				con.Ladder = string(byteLadder)
				con.IsMain = 1
				con.Desc = fmt.Sprintf("%v由于有流水自动生成", time.Now().Format("2006-01-02 15:04:05"))

				if _, err := models.AddContract(&con, order.GameID, 1, order.Cp, 1, ""); err != nil {
					beego.Error(err)
				}
				fmt.Printf("%v\n", con)
			}
		}
	}

}
