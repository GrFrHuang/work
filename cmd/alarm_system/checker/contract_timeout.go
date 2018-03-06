package checker

import (
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/cmd/alarm_system/alarm"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

// 渠道合同到期时间
// 有到期时间的渠道合同, 并且到期时间<当前时间+[半个月] (到期时间-当前时间<[半个月])
func ContractTimeout() (err error) {
	// 获取规则
	rules, err := models.GetAlarmRuleByType(models.AT_ContractTimeout)
	if err != nil {
		return
	}

	// 这里只有一个规则
	// 未定义报警规则
	if len(rules) == 0 {
		return
	}
	rule := rules[0]

	// 只用得到 时间
	var maxInterval int = 0
	maxInterval = rule.IntervalTime
	if rule.CustomIntervalTime != 0 {
		maxInterval = rule.CustomIntervalTime
	}
	rule = models.AlarmRule{IntervalTime: maxInterval}

	cs := []models.Contract{}
	_, err = orm.NewOrm().QueryTable(new(models.Contract)).
			Filter("company_type", 1).
			Filter("end_time__gt", 0).
			Filter("end_time__lt", int(time.Now().Unix()) + rule.IntervalTime).
			All(&cs)
	if err != nil {
		return
	}

	// 没有超时的
	if len(cs) == 0 {
		return
	}

	models.GroupRemitAddGameInfo(&cs)
	models.GroupRemitContractAddChannelInfo(&cs)
	models.GroupRemitAddCpInfo(&cs)
	models.GroupRemitAddContractStatusInfo(&cs)

	titles := []string{"渠道","公司名","合同签订时间","合同终止时间","我方主体","合作游戏"}
	contents := [][]string{}
	for _, n := range cs {
		channelName := ""
		if n.Channel != nil {
			channelName = n.Channel.Name
		}
		gameName := ""
		if n.Game != nil {
			gameName = n.Game.GameName
		}
		var body string
		if n.BodyMy == 1 {
			body = "云端"
		} else if n.BodyMy == 2 {
			body = "有量"
		}
		contents = append(contents,[]string{n.Channel.Name,n.ChannelCompanyName,
			time.Unix(n.SigningTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(n.EndTime, 0).Format("2006-01-02 15:04:05"),
			body,n.Game.GameName})

		// 发送邮件
		//SendContractTimeOut(n)

		// 得到还有几天过期
		days := (n.EndTime - time.Now().Unix()) / 3600 / 24
		err = alarm.SaveContractTimeout(channelName, gameName, days)
		if err != nil {
			return
		}
	}
	if len(contents)>0 {
		// 发送邮件
		SendTableMail("【渠道合同到期提醒】", titles, contents)
	}

	return
}
