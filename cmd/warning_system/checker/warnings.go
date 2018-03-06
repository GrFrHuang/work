package checker

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"strings"
)

// ChannelContractExpireWarning 渠道合同到期预警, 当渠道合同到期时按照相应的
// 预警规则发送邮件并添加预警日志
func ChannelContractExpireWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_CONTRACT_EXPIRE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}

	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	numbers := IntervalsToDates(intervals...)
	condition := ""
	if len(numbers) > 0 {
		holder := strings.Repeat(",?", len(numbers))
		condition = "HAVING FROM_UNIXTIME(max_end_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		return
	}

	var contracts []orm.Params
	_, err = orm.NewOrm().Raw(fmt.Sprintf("SELECT game_id, channel_code, MAX(end_time) AS max_end_time, begin_time, body_my "+
		"FROM contract "+
		"WHERE company_type=1  AND state not in (155,156) "+
		"GROUP BY game_id, channel_code  "+
		"%s", condition), numbers...).Values(&contracts)
	if err != nil || len(contracts) == 0 {
		return
	}

	var keys []string
	for _, contract := range contracts {
		gameId, _ := util.Interface2Int(contract["game_id"], false)
		endTime, _ := util.Interface2Int(contract["max_end_time"], false)
		channelCode, _ := util.Interface2String(contract["channel_code"], false)
		days := (endTime - time.Now().Unix()) / 3600 / 24
		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)
		keys = append(keys, key)
	}

	oldLog, _ := getRepeatWarningLog(rule.Type, keys)
	userEmails := GetEmailsByUserIds(rule.UserIds)

	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "渠道名称", "公司名称", "商务负责人", "游戏名称", "签订/终止时间"}
	var contents [][]string

	for _, contract := range contracts {
		gameId, _ := util.Interface2Int(contract["game_id"], false)
		endTime, _ := util.Interface2Int(contract["max_end_time"], false)
		channelCode, _ := util.Interface2String(contract["channel_code"], false)
		days := (endTime - time.Now().Unix()) / 3600 / 24
		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)

		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		beginTime, _ := util.Interface2Int(contract["begin_time"], false)
		channelCompany, _ := models.GetChannelCompanyByCode(channelCode)
		bodyId, _ := util.Interface2Int(contract["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))

		companyName := ""
		if channelCompany != nil {
			companyName = channelCompany.CompanyName
		}

		game, _ := models.GetGameByGameId(int(gameId))
		info := fmt.Sprintf("请注意,[%s]-[%s]的合同还有[%d]天到期",
			game.GameName,
			channelCompany.ChannelName,
			days)

		yunName := channelCompany.YunduanResPerName
		if yunName == "" {
			yunName = "无"
		}
		youName := channelCompany.YouliangResPerName
		if youName == "" {
			youName = "无"
		}

		contents = append(contents, []string{
			bodyMy,
			channelCompany.ChannelName,
			companyName,
			fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName),
			game.GameName,
			fmt.Sprintf("<span style='color: red'>%s~%s</span>", time.Unix(beginTime, 0).Format("2006-01-02"),
				time.Unix(int64(endTime), 0).Format("2006-01-02")),
		})

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      game.GameId,
			Grade:       rule.Grade,
			ChannelCode: channelCompany.ChannelCode,
			Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// CpContractExpireWarning CP合同到期预警, 当Cp合同到期时按照相应的
// 预警规则发送邮件并添加预警日志
func CpContractExpireWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CP_CONTRACT_EXPIRE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	numbers := IntervalsToDates(intervals...)

	condition := ""
	if len(numbers) > 0 {
		holder := strings.Repeat(",?", len(numbers))
		condition = "HAVING FROM_UNIXTIME(max_end_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		return
	}

	var contracts []orm.Params
	_, err = orm.NewOrm().Raw(fmt.Sprintf("SELECT game_id, company_id, MAX(end_time) AS max_end_time, begin_time, body_my "+
		"FROM contract "+
		"WHERE company_type=0 AND state not in (155,156) "+
		"GROUP BY game_id, channel_code  "+
		"%s", condition), numbers...).Values(&contracts)
	if err != nil || len(contracts) == 0 {
		return
	}

	var keys []string
	for _, contract := range contracts {
		gameId, _ := util.Interface2Int(contract["game_id"], false)
		endTime, _ := util.Interface2Int(contract["max_end_time"], false)
		companyId, _ := util.Interface2Int(contract["company_id"], false)
		days := (endTime - time.Now().Unix()) / 3600 / 24
		key := fmt.Sprintf("%d|%d|%d|%d", gameId, companyId, endTime, days)
		keys = append(keys, key)
	}

	oldLog, _ := getRepeatWarningLog(rule.Type, keys)
	userEmails := GetEmailsByUserIds(rule.UserIds)
	//userEmails := []string{"liuqilin@kuaifazs.com"}
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "CP名称", "商务负责人", "游戏名称", "签订/终止时间"}
	var contents [][]string

	for _, contract := range contracts {
		gameId, _ := util.Interface2Int(contract["game_id"], false)
		endTime, _ := util.Interface2Int(contract["max_end_time"], false)
		//companyId, _ := util.Interface2Int(contract["company_id"], false)
		days := (endTime - time.Now().Unix()) / 3600 / 24
		game, _ := models.GetGameByGameId(int(gameId))
		companyId := game.Issue

		key := fmt.Sprintf("%d|%d|%d|%d", gameId, companyId, endTime, days)

		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		beginTime, _ := util.Interface2Int(contract["begin_time"], false)
		company, _ := models.GetDistributionCompanyByCompanyId(companyId)
		bodyId, _ := util.Interface2Int(contract["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))


		info := fmt.Sprintf("请注意,[%s]-[%s]的合同还有[%d]天到期",
			game.GameName,
			company.CompanyName,
			days)

		yunName := company.YunduanResPerName
		if yunName == "" {
			yunName = "无"
		}
		youName := company.YouliangResPerName
		if youName == "" {
			youName = "无"
		}

		contents = append(contents, []string{
			bodyMy,
			company.CompanyName,
			fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName),
			game.GameName,
			fmt.Sprintf("<span style='color:red'>%s~%s</span>", time.Unix(beginTime, 0).Format("2006-01-02"),
				time.Unix(int64(endTime), 0).Format("2006-01-02")),
		})

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      game.GameId,
			Grade:       rule.Grade,
			//ChannelCode: strconv.Itoa(int(companyId)),
			Keys: key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// ChannelContractSignWarning 渠道合同未签订预警, 在游戏发布前
// 提前设定的天数进行预警
func ChannelContractSignWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_CONTRACT_SIGN)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = "AND FROM_UNIXTIME(publish_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}

	sql := fmt.Sprintf("SELECT g.`game_id`, g.`game_name`, ch.`cp`, g.`publish_time`, ch.`name`, c.`body_my` "+
		"FROM contract AS c LEFT JOIN game AS g  ON g.`game_id` = c.`game_id` LEFT JOIN channel AS ch ON c.channel_code=ch.cp "+
		"WHERE c.`company_type`=1 AND c.`state`=149 "+
		"%s "+
		"ORDER BY g.publish_time", condition)
	var contracts []orm.Params
	_, err = orm.NewOrm().Raw(sql, dates...).Values(&contracts)
	if err != nil || len(contracts) == 0 {
		return
	}

	oldLog := map[string]bool{}
	if rule.IsRepeat == 0 {
		var keys []string
		for _, contract := range contracts {
			gameId, _ := util.Interface2Int(contract["game_id"], false)
			publishTime, _ := util.Interface2Int(contract["publish_time"], false)
			channelCode, _ := util.Interface2String(contract["cp"], false)
			days := (publishTime - time.Now().Unix()) / 3600 / 24
			key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
			keys = append(keys, key)
		}
		oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "游戏名称", "渠道名称", "首发时间", "距离首发天数"}
	var contents [][]string
	for _, contract := range contracts {
		gameId, _ := util.Interface2Int(contract["game_id"], false)
		channelCode, _ := util.Interface2String(contract["cp"], false)
		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
		days := (publishTime - time.Now().Unix()) / 3600 / 24
		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		gameName, _ := util.Interface2String(contract["game_name"], false)
		channelName, _ := util.Interface2String(contract["name"], false)
		bodyId, _ := util.Interface2Int(contract["body_my"], false)

		bodyMy := ParseBodyMy(int(bodyId))
		contents = append(contents, []string{
			bodyMy,
			gameName,
			channelName,
			time.Unix(publishTime, 0).Format("2006-01-02"),
			//strconv.Itoa(int(days)),
			fmt.Sprintf("<span style='color:red'>%d</span>", days),
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]还有[%d]天首发，与渠道[%s]的合同还未签订",
			gameName,
			days,
			channelName)

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			ChannelCode: channelCode,
			Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// ChannelOrderVerifyWarning 渠道未对账预警，当流水达到指定值
// 时会在次月20日报警
func ChannelOrderVerifyWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_VERIFY)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	now := time.Now()
	lastMonth := now.Format("2006-01")
	if now.Day() != 20 {
		return
	}

	var data []orm.Params
	sql := "SELECT o.channel_code, con.body_my, ch.`name` AS channel_name, SUM(o.amount) AS total, o.`date` " +
		"FROM `order_pre_verify_channel` AS o LEFT JOIN channel AS ch ON o.`channel_code` = ch.`cp` " +
		"LEFT JOIN contract AS con ON o.`channel_code`=con.`channel_code` AND  con.company_type = 1 AND  o.game_id = con.game_id " +
		"AND con.effective_state=1 " +
		"WHERE o.verify_id = 0 AND o.`date` < ? " +
		"GROUP BY o.channel_code, con.body_my, o.`date` " +
		"HAVING total>? " +
		"ORDER BY o.`date`"
	_, err = orm.NewOrm().Raw(sql, lastMonth, rule.Amount).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}

	oldLog := map[string]bool{}
	if rule.IsRepeat == 0 {
		var keys []string
		for _, v := range data {
			date, _ := util.Interface2String(v["date"], false)
			channelCode, _ := util.Interface2String(v["channel_code"], false)
			key := fmt.Sprintf("%s|%s", channelCode, date)
			keys = append(keys, key)
		}
		oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "渠道名称", "我方流水", "未对账日期"}
	var contents [][]string
	for _, per := range data {
		date, _ := util.Interface2String(per["date"], false)
		channelCode, _ := util.Interface2String(per["channel_code"], false)
		key := fmt.Sprintf("%s|%s", channelCode, date)
		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		total, _ := util.Interface2Float(per["total"], false)
		channelName, _ := util.Interface2String(per["channel_name"], false)
		bodyId, _ := util.Interface2Int(per["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		contents = append(contents, []string{
			bodyMy,
			channelName,
			fmt.Sprintf("%.2f", total),
			fmt.Sprintf("<span style='color:red'>%s</span>", date),
		})

		info := fmt.Sprintf("请注意,渠道[%s]的流水在[%s]达到[%.2f]未对账", channelName, date, total)
		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			Grade:       rule.Grade,
			Amount:      total,
			Date:        date,
			ChannelCode: channelCode,
			Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// CpOrderVerifyWarning CP未对账预警，当流水达到指定值
// 时会在次月20日报警
func CpOrderVerifyWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CP_ORDER_VERIFY)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	now := time.Now()
	month := now.Format("2006-01")
	if now.Day() != 20 {
		return
	}

	var data []orm.Params
	sql := "SELECT o.company_id, con.body_my, ch.`name` AS company_name, SUM(o.amount) AS total, o.`date` " +
		"FROM `order_pre_verify_cp` AS o LEFT JOIN company_type AS ch ON o.`company_id` = ch.`id` " +
		"LEFT JOIN contract AS con ON o.company_id = con.company_id AND  con.company_type = 0 AND  o.game_id = con.game_id " +
		"WHERE o.verify_id = 0 AND o.`date` < ? " +
		"GROUP BY o.company_id, con.body_my, o.`date` " +
		"HAVING total>? " +
		"ORDER BY o.`date`"
	_, err = orm.NewOrm().Raw(sql, month, rule.Amount).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}

	oldLog := map[string]bool{}
	if rule.IsRepeat == 0 {
		var keys []string
		for _, v := range data {
			date, _ := util.Interface2String(v["date"], false)
			companyId, _ := util.Interface2String(v["company_id"], false)
			key := fmt.Sprintf("%s|%s", companyId, date)
			keys = append(keys, key)
		}
		oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "发行商", "我方流水", "未对账日期"}
	var contents [][]string
	for _, per := range data {
		date, _ := util.Interface2String(per["date"], false)
		companyId, _ := util.Interface2String(per["company_id"], false)
		key := fmt.Sprintf("%s|%s", companyId, date)
		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		total, _ := util.Interface2Float(per["total"], false)
		companyName, _ := util.Interface2String(per["company_name"], false)
		bodyId, _ := util.Interface2Int(per["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		contents = append(contents, []string{
			bodyMy,
			companyName,
			fmt.Sprintf("%.2f", total),
			fmt.Sprintf("<span style='color:red'>%s</span>", date),
		})

		info := fmt.Sprintf("请注意,发行商[%s]的流水在[%s]达到[%.2f]未对账", companyName, date, total)
		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			Grade:       rule.Grade,
			Amount:      total,
			Date:        date,
			//ChannelCode: companyId,
			Keys: key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// ChannelOrderThreshold 流水总额预警，当渠道的某个游戏流水达到指定的阈值时
// 会发送预警日志和插入预警日志
func ChannelOrderThresholdWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_THRESHOLD)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	var data []orm.Params
	sql := "SELECT o.game_id, o.cp AS channel_code, g.`game_name`, c.`name` AS channel_name, o.total, con.body_my " +
		"FROM( SELECT cp, game_id, SUM(amount) AS total FROM `order` " +
		"GROUP BY game_id, cp) AS o LEFT JOIN `game_all` AS g ON o.`game_id` = g.`game_id` LEFT JOIN channel AS c ON o.cp = c.`cp` LEFT JOIN contract AS con ON o.`cp` = con.`channel_code` AND o.`game_id`=con.`game_id` " +
		"WHERE o.total > ? " +
		"ORDER BY o.total DESC "
	_, err = orm.NewOrm().Raw(sql, rule.Amount).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}

	oldLog := map[string]bool{}
	if rule.IsRepeat == 0 {
		var keys []string
		for _, v := range data {
			channelCode, _ := util.Interface2String(v["channel_code"], false)
			gameId, _ := util.Interface2String(v["game_id"], false)
			key := fmt.Sprintf("%s|%s", channelCode, gameId)
			keys = append(keys, key)
		}
		oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "渠道商", "游戏名称", "流水"}
	var contents [][]string
	for _, per := range data {
		channelCode, _ := util.Interface2String(per["channel_code"], false)
		gameId, _ := util.Interface2Int(per["game_id"], false)
		key := fmt.Sprintf("%s|%d", channelCode, gameId)
		if rule.IsRepeat == 0 {
			if _, ok := oldLog[key]; ok {
				continue
			}
		}

		total, _ := util.Interface2Float(per["total"], false)
		channelName, _ := util.Interface2String(per["channel_name"], false)
		gameName, _ := util.Interface2String(per["game_name"], false)
		bodyId, _ := util.Interface2Int(per["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))

		contents = append(contents, []string{
			bodyMy,
			channelName,
			gameName,
			fmt.Sprintf("<span style='color:red'>%.2f</span>", total),
		})

		info := fmt.Sprintf("请注意,[%s]下的[%s]流水达到[%.2f]，已经超过[%.2f]阈值", channelName, gameName, total, rule.Amount)
		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			Grade:       rule.Grade,
			Amount:      total,
			ChannelCode: channelCode,
			GameId:      int(gameId),
			Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// ChannelOrderPayWarning 渠道未回款预警，当渠道的未回款金额达到指定的阈值时
// 会发送预警日志和插入预警日志
func ChannelOrderPayWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_PAY)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	var data []orm.Params
	sql := "SELECT c.body_my, c.remit_company_id, d.name, IFNULL(a.amount, 0) - IFNULL(b.amount, 0) AS total_amount " +
		"FROM remit_pre_account AS c " +
		"LEFT JOIN ( " +
		"SELECT body_my,remit_company_id,SUM(amount_payable) AS amount " +
		"FROM `verify_channel` WHERE `status` >20 " +
		"GROUP BY body_my,remit_company_id) AS a ON c.remit_company_id = a.remit_company_id AND c.body_my = a.body_my " +
		"LEFT JOIN ( " +
		"SELECT body_my,remit_company_id,SUM(amount) AS amount " +
		"FROM `remit_down_account` " +
		"GROUP BY body_my,remit_company_id) AS b ON c.remit_company_id = b.remit_company_id AND c.body_my = b.body_my " +
		"LEFT JOIN company_type AS d ON c.remit_company_id = d.id " +
		"HAVING total_amount > ? " +
		"ORDER BY total_amount DESC"
	_, err = orm.NewOrm().Raw(sql, rule.Amount).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}

	// 没有办法进行是否为重复消息
	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)

	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "回款主体", "渠道名称(商务负责人)", "未回款金额"}
	var contents [][]string
	for _, per := range data {
		//channelCode, _ := util.Interface2String(per["body_my"], false)
		//gameId, _ := util.Interface2String(per["game_id"], false)
		//key := fmt.Sprintf("%s|%s", channelCode, gameId)
		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}

		total, _ := util.Interface2Float(per["total_amount"], false)
		companyName, _ := util.Interface2String(per["name"], false)
		companyId, _ := util.Interface2Int(per["remit_company_id"], false)
		bodyId, _ := util.Interface2Int(per["body_my"], false)
		//company_id, _ := util.Interface2String(per["remit_company_id"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		channelCodes := models.GetChannelAndPeople(companyId)

		contents = append(contents, []string{
			bodyMy,
			companyName,
			channelCodes,
			fmt.Sprintf("<span style='color:red'>%.2f</span>", total),
		})

		info := fmt.Sprintf("请注意,[%s]与[%s]未回款金额达到[%.2f]，已经超过[%.2f]阈值", companyName, bodyMy, total, rule.Amount)
		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			Grade:       rule.Grade,
			//ChannelCode: company_id,
			Amount: total,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}
