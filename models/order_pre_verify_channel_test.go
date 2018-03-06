package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	"github.com/tealeg/xlsx"
	"time"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"encoding/json"
)

// 导出渠道未对账
func TestExportNotVerify(t *testing.T) {
	_, notVerify, _ := GetChannelNotVerifyInfo(1000, 0)

	var rs []map[string]interface{}
	for _, verify := range notVerify {
		data, _ := GetNotVerifyChannelGame(verify.BodyMy, verify.Channel.Cp, verify.Date)
		for _, dat := range data {
			r := map[string]interface{}{}

			r["date"] = verify.Date
			r["channel"] = verify.Channel.Name
			if verify.BodyMy == 1 {
				r["body"] = "云端"
			} else if verify.BodyMy == 2 {
				r["body"] = "有量"
			}
			r["game"] = dat.GameName
			r["amount_my"] = dat.Amount
			r["amount_theory"] = dat.AmountTheory
			r["amount_opposite"] = dat.AmountOpposite
			r["amount_payable"] = dat.AmountPayable
			r["difference"] = dat.AmountTheory - dat.AmountPayable

			rs = append(rs, r)
		}
	}

	cols := []string{"date", "channel", "body", "game", "amount_my", "amount_theory", "amount_opposite", "amount_payable", "difference"}
	maps := map[string]string{
		cols[0]: "对账日期",
		cols[1]: "渠道名称",
		cols[2]: "我方主体",
		cols[3]: "游戏名称",
		cols[4]: "我方流水",
		cols[5]: "理论金额",
		cols[6]: "对方流水",
		cols[7]: "应收金额",
		cols[8]: "差额",
	}

	file := xlsx.NewFile()
	sht, err := file.AddSheet("渠道未对账")
	if err != nil {
		return
	}
	row := sht.AddRow()
	for _, c := range cols {
		if maps != nil {
			row.AddCell().Value = maps[c]
		} else {
			row.AddCell().Value = c
		}
	}
	// 添加一个空行
	sht.AddRow()

	for _, row := range rs {
		newRow := sht.AddRow()
		for _, value := range cols {
			if v, ok := row[value]; ok {
				newRow.AddCell().Value = fmt.Sprintf("%v", v)
			} else {
				newRow.AddCell().Value = ""
			}
		}
	}

	err = file.Save("C:/Users/Administrator/Desktop/渠道未对账2.xlsx")
	if err != nil {
		return
	}
}

// 根据渠道对账导入
func TestImportVerify(t *testing.T) {
	excelFileName := "C:/Users/Administrator/Desktop/渠道未对账1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet3" {
			for _, row := range sheet.Rows {
				date, _ := row.Cells[0].String()
				channel_code, _ := row.Cells[1].String()
				body_my, _ := row.Cells[2].Int()
				game_id, _ := row.Cells[3].String()
				amount, _ := row.Cells[4].Float()
				amount_theory, _ := row.Cells[5].Float()
				amount_opposite, _ := row.Cells[6].Float()
				amount_payable, _ := row.Cells[7].Float()

				flag := orm.NewOrm().QueryTable(new(VerifyChannel)).Filter("date", date).Filter("body_my", body_my).
					Filter("channel_code", channel_code).Exist()

				if flag == true {
					verify := VerifyChannel{}
					orm.NewOrm().QueryTable(new(VerifyChannel)).Filter("date", date).Filter("body_my", body_my).
						Filter("channel_code", channel_code).One(&verify)
					verify.AmountMy = verify.AmountMy + amount
					verify.AmountOpposite = verify.AmountOpposite + amount_opposite
					verify.AmountPayable = verify.AmountPayable + amount_payable
					verify.AmountTheory = verify.AmountTheory + amount_theory
					//orm.NewOrm().Update(&verify)
					orm.NewOrm().Raw("UPDATE order_pre_verify_channel SET amount_opposite=?,amount_payable= ?,verify_id=? "+
						"WHERE game_id=? AND DATE =? AND channel_code= ?;", amount_opposite, amount_payable, verify.Id,
						game_id, date, channel_code).Exec()
				} else {
					verify := VerifyChannel{
						Date:           date,
						BodyMy:         body_my,
						ChannelCode:    channel_code,
						AmountMy:       amount,
						AmountOpposite: amount_opposite,
						AmountPayable:  amount_payable,
						AmountTheory:   amount_theory,
						Status:         10,
						VerifyTime:     int(time.Now().Unix()),
						VerifyUserId:   13,
						FileId:         0,
						FilePreviewId:  0,
						Desc:           "2018-02-08手动导入",
						CreatedTime:    int(time.Now().Unix()),
						CreatedUserId:  13,
						UpdatedTime:    int(time.Now().Unix()),
						UpdatedUserId:  13,
					}
					id, _ := orm.NewOrm().Insert(&verify)
					orm.NewOrm().Raw("UPDATE order_pre_verify_channel SET amount_opposite=?,amount_payable= ?,verify_id=? "+
						"WHERE game_id=? AND DATE =? AND channel_code= ?;", amount_opposite, amount_payable, id,
						game_id, date, channel_code).Exec()
					fmt.Printf("--------------")
				}
			}
		}
	}
}

func TestImportChannelAccessAndContract(t *testing.T) {
	excelFileName := "C:/Users/Administrator/Desktop/接入记录l(1) (1).xlsx"
	// 游戏名	渠道	我方主体	合作方式	我方比例(小数)	通道费(小数)	接入人	接入时间	合同开始时间	合同终止时间	合同签订状态
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet1" {
			for _, row := range sheet.Rows {
				gameName, _ := row.Cells[0].String()
				channelName, _ := row.Cells[1].String()
				bodyMy, _ := row.Cells[2].String()
				my, _ := row.Cells[4].Float()
				slot, _ := row.Cells[5].Float()
				userName, _ := row.Cells[6].String()
				//importTime, _ := row.Cells[7].GetTime(false)
				startTime, _ := row.Cells[8].GetTime(false)
				endTime, _ := row.Cells[9].GetTime(false)

				var v ChannelAccess
				game := GetGameByGamename(gameName)
				v.GameId = game.GameId
				v.PublishTime = game.PublishTime
				channel, _ := GetChannelByChannelName(channelName)
				channelCompany, _ := GetChannelCompanyByCode(channel.Cp)
				v.ChannelCode = channel.Cp
				if bodyMy == "云端" {
					v.BodyMy = 1
					v.BusinessPerson = channelCompany.YunduanResponsiblePerson
				} else {
					v.BodyMy = 2
					v.BusinessPerson = channelCompany.YouliangResponsiblePerson
				}
				v.Cooperation = 45
				user, _ := GetUserInfoByNickName(userName)
				v.UpdateChannelUserID = user.Id
				v.UpdateChannelTime = time.Now().Unix()

				v.AccessState = 1
				v.AccessUpdateUser = user.Id
				v.AccessUpdateTime = time.Now().Unix()

				ladder := []old_sys.Ladder4Post{{
					SlottingFee: slot,
					Ratio:       my,
				}}
				byteLadder, _ := json.Marshal(ladder)
				v.Ladders = string(byteLadder)

				if _, err1 := AddChannelAccess(&v); err1 != nil {
					fmt.Printf("------err:%v\n", err1)
				}
				var con Contract
				var conn GamePlanChannel
				con.Ladder = v.Ladders
				con.BodyMy = v.BodyMy
				con.IsMain = 1
				con.Desc = "2018-02-22手动导入"
				con.BeginTime = startTime.Unix()
				con.EndTime = endTime.Unix()
				// 如果渠道为“预充值”，则该合同为“预充值渠道，无需合同”
				if v.Cooperation == 166 {
					con.State = 164
				}

				if _, err := AddContract(&con, v.GameId, 1, v.ChannelCode, user.Id, v.Accessory); err != nil {
					fmt.Printf("------err2:%v\n", err)
				}
				if _, err := AddGamePlanChannel(&conn, v.GameId, v.ChannelCode); err != nil {
					fmt.Printf("------err3:%v\n", err)
				}
			}
		}
	}
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
		"kuaifazs", "10.8.230.17",
		"3308", "work_together_online")
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
	//	"123456", "10.8.230.17",
	//	"3308", "work_together")
	fmt.Printf("link:%v", link)
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
