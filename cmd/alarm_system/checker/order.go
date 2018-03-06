package checker

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/log"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

// todo 整个流水都没做, 由于不知道流水计算的开始时间与截止时间
// 某游戏某渠道的流水报警
func Order() (err error) {
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

	not, err := GetAllNotVerify()
	if err != nil {
		log.Error("alarm-remit", err)
		return
	}
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

		//log.Verbose("alarm-verify", n)
		if n.Amount > rule.Amount || int(time.Now().Unix())-n.StartTime > rule.IntervalTime {
			//err = alarm.SaveVerifyAlarmLog(n.ChannelCode, n.StartTime, "", n.Amount)
			//if err != nil {
			//	return
			//}
		}
	}

	return
}

//
//type NotVerify struct {
//	ChannelCode string
//	ChannelName string
//	Amount      float64
//	StartTime   int
//}

// 获取某游戏某渠道的流水
func GetAllAmount(channelCodes []string, games []int) (amount float64, err error) {
	return
}
