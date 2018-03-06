package checker

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"kuaifa.com/kuaifa/work-together/cmd/alarm_system/alarm"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"time"
)

type NotRemit struct {
	RemitCompanyId   int
	RemitCompanyName string
	Amount           float64
	StartTime        int64
}

// 检查未回款的回款主体
func RemitAccount() (err error) {
	// 获取规则
	rules, err := models.GetAlarmRuleByType(models.AT_RemitAccount)
	if err != nil {
		return
	}

	ruleMap := map[int]models.AlarmRule{} // RemitCompanyId => Rule
	for _, rule := range rules {
		var maxAmount float64 = 0
		var maxInterval int = 0

		if rule.Readonly == 1 {
			// 系统里自带的
			// 如果有自定义的值就使用自定义的值
			maxAmount = rule.Amount
			if rule.CustomAmount != 0 {
				maxAmount = rule.CustomAmount
			}
			maxInterval = rule.IntervalTime
			if rule.CustomIntervalTime != 0 {
				maxInterval = rule.CustomIntervalTime
			}
			// 默认值得key为0,如果没找到自定义回款主体的就按这个默认值
			ruleMap[0] = models.AlarmRule{Amount: maxAmount, IntervalTime: maxInterval}
		} else {
			// 自定义的
			maxAmount = rule.Amount
			maxInterval = rule.IntervalTime
			ruleMap[rule.RemitCompanyId] = models.AlarmRule{Amount: maxAmount, IntervalTime: maxInterval}
		}
	}

	// 未定义报警规则
	if len(ruleMap) == 0 {
		return
	}

	notRemit, err := GetAllNotRemitAccount()
	if err != nil {
		log.Error("alarm-remit", err)
		return
	}
	titles := []string{"公司名","未对账金额","未对账时间"}
	contents := [][]string{}
	for _, remit := range notRemit {
		// 有自定义的规则
		rule, ok := ruleMap[remit.RemitCompanyId]
		if !ok {
			rule, ok = ruleMap[0]
			// 没有默认规则
			if !ok {
				return
			}
		}

		//log.Verbose("alarm-remit",remit,int(time.Now().Unix())-remit.StartTime)
		if remit.Amount > rule.Amount || int(time.Now().Unix()-remit.StartTime) > rule.IntervalTime {
			contents = append(contents,[]string{remit.RemitCompanyName,
				strconv.FormatFloat(remit.Amount,'f',2,64),time.Unix(int64(remit.StartTime), 0).Format("2006-01-02")})

			//SendRemitAccount(remit)
			err = alarm.SaveRemitAlarmLog(remit.RemitCompanyId, int(remit.StartTime), remit.Amount)
			if err != nil {
				return
			}
		}
	}
	if len(contents)>0 {
		// 发送邮件
		SendTableMail("【渠道未回款提醒】", titles, contents)
	}

	return
}

// 所有未回款的渠道
// 获取所有渠道对账单中待回款的对账单的渠道
func GetAllNotRemitAccount() (notRemit []NotRemit, err error) {
	not, err := getAllChannelPreRemitAmount()
	if err != nil {
		return
	}

	AddCompanyName(&not)

	notRemit = []NotRemit{}
	for _, p := range not {
		amount, _ := strconv.ParseFloat(p["amount"].(string), 64)
		companyName := p["CompanyName"].(string)
		remitCompanyId, _ := strconv.Atoi(p["remit_company_id"].(string))
		t, _ := time.Parse("2006-01-02", p["start_time"].(string))
		firstTime := t.Unix()
		fmt.Println("firstTime:", firstTime)
		notRemit = append(notRemit, NotRemit{
			RemitCompanyId:   remitCompanyId,
			RemitCompanyName: companyName,
			Amount:           amount,
			StartTime:        firstTime,
		})
	}
	return
}

// 获取所有渠道,未回款金额
// {amount,remit_company_id,end_time,start_time}
func getAllChannelPreRemitAmount() (maps []orm.Params, err error) {
	maps = []orm.Params{}
	finishStatus := models.CHAN_VERIFY_S_RECEIPT
	sql := fmt.Sprintf("SELECT remit_company_id, max(end_time) as end_time,min(start_time) as start_time, SUM(amount_payable - amount_remit) as amount from channel_verify_account WHERE amount_payable != amount_remit AND status = %d GROUP By remit_company_id", finishStatus)
	o := orm.NewOrm()
	_, err = o.Raw(sql).Values(&maps)
	if err != nil {
		return
	}

	return
}

func AddCompanyName(maps *[]orm.Params) {
	if maps == nil || len(*maps) == 0 {
		return
	}

	companyIds := make([]int, len(*maps))
	for i, v := range *maps {
		companyIds[i], _ = strconv.Atoi(v["remit_company_id"].(string))
	}

	companies := []models.Company{}
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.Company)).
		Filter("id__in", companyIds).
		All(&companies, "id", "Name")
	if err != nil {
		return
	}

	linkMap := map[int]models.Company{}
	for _, g := range companies {
		linkMap[g.Id] = g
	}

	for i, s := range *maps {
		company_id, _ := strconv.Atoi(s["remit_company_id"].(string))
		if g, ok := linkMap[company_id]; ok {
			(*maps)[i]["CompanyName"] = g.Name
		}
	}
	return
}
