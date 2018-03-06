package checker

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"time"
)

// 报警 只需要报渠道对账
func VerifyAccount() (err error) {
	// 获取规则
	rules, err := models.GetAlarmRuleByType(models.AT_VerifyAccount)
	if err != nil {
		return
	}

	ruleMap := map[string]models.AlarmRule{} // ChannelCode => Rule
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
			ruleMap[""] = models.AlarmRule{Amount: maxAmount, IntervalTime: maxInterval}
		} else {
			// 自定义的
			maxAmount = rule.Amount
			maxInterval = rule.IntervalTime
			r := models.AlarmRule{Amount: maxAmount, IntervalTime: maxInterval}
			chans := []string{}
			json.Unmarshal([]byte(rule.ChannelCodes), &chans)
			for _, c := range chans {
				ruleMap[c] = r
			}
		}
	}

	// 未定义报警规则
	if len(ruleMap) == 0 {
		return
	}

	not, err := GetAllNotVerifyChannel()
	if err != nil {
		log.Error("alarm-verfiy", err)
		return
	}
	titles := []string{"渠道", "我方主体", "未对账流水", "最后一次未对账"}
	contents := [][]string{}
	for _, n := range not {
		// 有自定义的规则
		rule, ok := ruleMap[n.ChannelCode]
		if !ok {
			rule, ok = ruleMap[""]
			// 没有默认规则
			if !ok {
				return
			}
		}

		isTimeout := false
		t, e := time.Parse("2006-01", n.Date)
		if e != nil {
			err = e
			return
		}
		isTimeout = int(time.Now().Unix()-t.Unix()) > rule.IntervalTime

		if n.Amount > rule.Amount || isTimeout {
			contents = append(contents,
				[]string{
					n.ChannelName,
					n.BodyMy,
					strconv.FormatFloat(n.Amount, 'f', 2, 64),
					n.Date,
				})

			//err = alarm.SaveVerifyAlarmLog(n.ChannelCode, n.Date, "", n.Amount)
			//if err != nil {
			//	return
			//}
		}
	}
	if len(contents) > 0 {
		// 发送邮件
		SendTableMail("【渠道未对账提醒】", titles, contents)
	}

	return
}

type NotVerify struct {
	ChannelCode string
	ChannelName string
	CompanyName string
	Amount      float64
	StartTime   int
}

type NotVerifyCp struct {
	CompanyName string  // 发行商名称
	Amount      float64 // 金额
	Date        string  // 月份
}

type NotVerifyChannel struct {
	CompanyName string  // 公司名称
	ChannelName string  // 渠道名
	BodyMy      string  // 云端 or 有量 or 未知 (未知是由于 游戏与渠道的合同上没有body_my为0)
	ChannelCode string  // 渠道code
	Amount      float64 // 金额
	Date        string  // 月份
}

// unused, 由于cp对账不用报警
// 获取所有cp未对账的 {发行商,金额,日期}
func GetAllNotVerifyCp() (notVerify []NotVerifyCp, err error) {
	sql := "SELECT company_id,company.`name` as company_name,sum(amount) as amount,min(date) as date " +
		"FROM `order_pre_verify_cp` " +
		"LEFT JOIN company_type ON company.id = company_id " +
		"WHERE verify_id = 0 " +
		"GROUP BY company_id"
	maps := []orm.Params{}
	_, err = orm.NewOrm().Raw(sql).Values(&maps)
	if err != nil {
		return
	}
	notVerify = []NotVerifyCp{}
	for _, v := range maps {
		companyName, _ := util.Interface2String(v["company_name"], false)
		// 这里由于有些游戏,没有填写发行商,所以导致的预对账单也没有发行商
		if companyName != "" {
			amount, _ := util.Interface2Float(v["amount"], false)
			date, _ := util.Interface2String(v["date"], false)
			notVerify = append(notVerify, NotVerifyCp{
				Amount:      amount,
				CompanyName: companyName,
				Date:        date,
			})
		}
	}

	return
}

// 获取所有cp未对账的 {渠道名,金额,日期}
func GetAllNotVerifyChannel() (notVerify []NotVerifyChannel, err error) {
	nowMonth := time.Now().Format("2006-01")
	sql := "SELECT channel_code,channel.`name` as channel_name,sum(amount) as amount,min(date) as date,contract.body_my " +
		"FROM `order_pre_verify_channel` " +
		"LEFT JOIN channel ON channel.cp = channel_code " +
		"LEFT JOIN contract ON contract.company_id = channel.channel_id AND  contract.company_type = 1 AND contract.game_id = order_pre_verify_channel.game_id " +
		"WHERE verify_id = 0 AND date != '" + nowMonth + "' " +
		"GROUP BY channel_code,body_my"
	maps := []orm.Params{}
	_, err = orm.NewOrm().Raw(sql).Values(&maps)
	if err != nil {
		return
	}
	notVerify = []NotVerifyChannel{}
	for _, v := range maps {
		channelName, _ := util.Interface2String(v["channel_name"], false)
		amount, _ := util.Interface2Float(v["amount"], false)
		date, _ := util.Interface2String(v["date"], false)
		channelCode, _ := util.Interface2String(v["channel_code"], false)
		bodyMy, _ := util.Interface2Int(v["body_my"], false)
		bodyMyStr := "有量"
		switch bodyMy {
		case 0:
			bodyMyStr = "未知"
		case 1:
			bodyMyStr = "云端"
		case 2:
			bodyMyStr = "有量"
		}
		notVerify = append(notVerify, NotVerifyChannel{
			Amount:      amount,
			ChannelName: channelName,
			ChannelCode: channelCode,
			BodyMy:      bodyMyStr,
			Date:        date,
		})
	}

	return
}

// 获取未对账的order
func GetAllNotVerify() (notVerify []NotVerify, err error) {
	maps := []orm.Params{}
	sql := fmt.Sprintf("SELECT " +
		"SUM(amount) as total,MIN(date) as start_time,cp " +
		"FROM `order` " +
		"WHERE channel_verified != 1 " +
		"GROUP BY cp")
	o := orm.NewOrm()
	_, err = o.Raw(sql).Values(&maps)
	if err != nil {
		return
	}

	addChannelName(&maps)
	addChannelCompanyId(&maps)
	addChannelCompanyName(&maps)

	notVerify = []NotVerify{}
	for _, p := range maps {
		amount, _ := strconv.ParseFloat(p["total"].(string), 64)
		channelCode := p["cp"].(string)
		companyName := ""
		if p["company_name"] != nil {
			companyName = p["company_name"].(string)
		}
		channelName := ""
		if p["channel_name"] != nil {
			channelName = p["channel_name"].(string)
		}
		t, _ := time.Parse("2006-01-02", p["start_time"].(string))
		firstTime := int(t.Unix())

		notVerify = append(notVerify, NotVerify{
			ChannelCode: channelCode,
			ChannelName: channelName,
			CompanyName: companyName,
			Amount:      amount,
			StartTime:   firstTime,
		})
	}

	return
}

func addChannelName(maps *[]orm.Params) {
	if maps == nil || len(*maps) == 0 {
		return
	}

	channel_codes := make([]string, len(*maps))
	for i, v := range *maps {
		channel_codes[i], _ = v["cp"].(string)
	}

	o := orm.NewOrm()
	channels := []models.Channel{}
	_, err := o.QueryTable(new(models.Channel)).
		Filter("cp__in", channel_codes).
		All(&channels, "cp", "name")
	if err != nil {
		return
	}

	linkMap := map[string]models.Channel{}
	for _, g := range channels {
		linkMap[g.Cp] = g
	}

	for i, s := range *maps {
		channel_code, _ := s["cp"].(string)
		if g, ok := linkMap[channel_code]; ok {
			(*maps)[i]["channel_name"] = g.Name
		}
	}
	return
}

func addChannelCompanyId(maps *[]orm.Params) {
	if maps == nil || len(*maps) == 0 {
		return
	}

	channel_codes := make([]string, len(*maps))
	for i, v := range *maps {
		channel_codes[i], _ = v["cp"].(string)
	}

	o := orm.NewOrm()
	channel_companies := []models.ChannelCompany{}
	_, err := o.QueryTable(new(models.ChannelCompany)).
		Filter("channel_code__in", channel_codes).
		All(&channel_companies, "channel_code", "company_id")
	if err != nil {
		return
	}

	linkMap := map[string]models.ChannelCompany{}
	for _, g := range channel_companies {
		linkMap[g.ChannelCode] = g
	}

	for i, s := range *maps {
		channel_code, _ := s["cp"].(string)
		if g, ok := linkMap[channel_code]; ok {
			(*maps)[i]["company_id"] = g.CompanyId
		}
	}
	return
}

func addChannelCompanyName(maps *[]orm.Params) {
	if maps == nil || len(*maps) == 0 {
		return
	}

	companyIds := make([]int, len(*maps))
	for i, v := range *maps {
		if v["company_id"] != nil {
			companyIds[i], _ = v["company_id"].(int)
		} else {
			companyIds[i] = 45
		}

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
		if s["company_id"] != nil {
			tmp_companyId := s["company_id"].(int)
			if g, ok := linkMap[tmp_companyId]; ok {
				(*maps)[i]["company_name"] = g.Name
			}
		}
	}
	return
}
