package checker

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/cmd/alarm_system/alarm"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

// 渠道合同在游戏上线前半个月还未签订
func ContractNotSign() (err error) {
	// 获取规则
	rules, err := models.GetAlarmRuleByType(models.AT_ContractSign)
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

	csMap := []orm.Params{}
	_, err = orm.NewOrm().Raw("SELECT contract.game_id,company.name as company_name,channel.name as channel_name,game.game_name,game.publish_time " +
			"FROM `contract` LEFT JOIN `channel` on channel.channel_id = contract.company_id RIGHT JOIN `game` on contract.game_id = game.game_id LEFT JOIN `company` on contract.company_id = company.id WHERE contract.company_type = 1 AND contract.state=0",//1 AND game.publish_time>0 AND game.publish_time<?",
		).Values(&csMap)
	if err != nil {
		return
	}
	titles := []string{"公司|渠道","合作游戏"}
	contents := [][]string{}

	// channelName=>games
	mailData := map[string]string{}
	for _, c := range csMap {
		gameName, _ := util.Interface2String(c["game_name"], true)
		companyName, _ := util.Interface2String(c["company_name"], true)
		publishTime, _ := util.Interface2Int(c["publish_time"], false)
		channelName, _ := util.Interface2String(c["channel_name"], true)

		key := companyName + "," + channelName
		if _, ok := mailData[key]; ok {
			mailData[key] = mailData[key] + "," + gameName
		} else {
			mailData[key] = gameName
		}

		// 得到还有几天发行
		days := (publishTime - time.Now().Unix()) / 3600 / 24
		err = alarm.SaveContractNotSign(channelName, gameName, days)
		if err != nil {
			return
		}
	}

	for channelName, gameName := range mailData {
		contents = append(contents,[]string{channelName,gameName})

		//SendContractNotSign(channelName, gameName)
	}

	if len(contents)>0 {
		// 发送邮件
		SendTableMail("【渠道合同未签订提醒】", titles, contents)
	}

	return
}
