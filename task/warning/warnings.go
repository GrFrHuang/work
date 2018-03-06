package warning

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"strings"
	"strconv"
	"kuaifa.com/kuaifa/work-together/cmd/warning_system/checker"
)

// 游戏停运公告
func GameOutageWarning(gameid ...int) (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_OUTAGE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(true, intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = "WHERE FROM_UNIXTIME(incr_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:]) +
			" OR FROM_UNIXTIME(recharge_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:]) +
			" OR FROM_UNIXTIME(server_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}
	if len(gameid) > 0 {
		condition = fmt.Sprintf(" WHERE a.game_id=%v", gameid[0])
	}

	sql := fmt.Sprintf("SELECT a.game_id,b.game_name,a.incr_time,a.recharge_time,a.server_time,a.desc,a.create_person FROM game_outage a "+
		"LEFT JOIN game b ON a.game_id=b.game_id "+
		"%s "+
		"ORDER BY a.server_time desc ", condition)
	var gameOutages []orm.Params

	if len(gameid) > 0 {
		_, err = orm.NewOrm().Raw(sql).Values(&gameOutages)
	} else {
		_, err = orm.NewOrm().Raw(sql, dates, dates, dates).Values(&gameOutages)
	}
	if err != nil || len(gameOutages) == 0 {
		return
	}

	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//	keys := []string{}
	//	for _, contract := range gameOutages {
	//		gameId, _ := util.Interface2Int(contract["game_id"], false)
	//		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
	//		channelCode, _ := util.Interface2String(contract["cp"], false)
	//		days := (publishTime - time.Now().Unix()) / 3600 / 24
	//		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
	//		keys = append(keys, key)
	//	}
	//	oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"游戏名称", "添加人", "关闭新增时间", "关闭充值时间", "关闭服务器时间", "备注"}
	var contents [][]string
	for _, outage := range gameOutages {
		gameId, _ := util.Interface2Int(outage["game_id"], false)
		createPerson, _ := util.Interface2Int(outage["create_person"], false)
		incrTime, _ := util.Interface2Int(outage["incr_time"], false)
		rechargeTime, _ := util.Interface2Int(outage["recharge_time"], false)
		serverTime, _ := util.Interface2Int(outage["server_time"], false)
		desc, _ := util.Interface2String(outage["desc"], false)

		incrDays := TimeSubDays(time.Unix(incrTime, incrTime), time.Now())
		rechargeDays := TimeSubDays(time.Unix(rechargeTime, rechargeTime), time.Now())
		serverDays := TimeSubDays(time.Unix(serverTime, serverTime), time.Now())
		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}

		per, _ := strconv.Atoi(strconv.FormatInt(createPerson, 10))
		person, _ := models.GetUserById(per)
		gameName, _ := util.Interface2String(outage["game_name"], false)

		contents = append(contents, []string{
			gameName,
			person.Nickname,
			time.Unix(incrTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(rechargeTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(serverTime, 0).Format("2006-01-02 15:04:05"),
			desc,
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]即将停运!", gameName)
		if incrDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止新增!", incrDays)
		}
		if rechargeDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止充值!", rechargeDays)
		}
		if serverDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将关闭服务器!", serverDays)
		}

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			//ChannelCode: channelCode,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// 下架游戏合同
func DownContract() (err error) {
	timeNow := time.Now().Unix()
	sql := fmt.Sprintf("SELECT b.* "+
		"FROM game_outage a LEFT JOIN contract b ON a.game_id=b.game_id "+
		"WHERE a.server_time<%v AND b.state != 155 ", timeNow)

	o := orm.NewOrm()
	var contracts []models.Contract

	_, err = o.Raw(sql).QueryRows(&contracts)

	if len(contracts) == 0 {
		return
	}
	for _, con := range contracts {
		// 将达到关服时间的所有游戏的合同状态改为无合作
		con.State = 155
		_, err = o.Update(&con)
	}

	return
}

// 游戏停运渠道合同待处理
func GameOutageChannelContractOpe() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_OUTAGE_CHANNEL_OPERATE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(true, intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = " AND FROM_UNIXTIME(ou.server_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}

	sql := fmt.Sprintf("SELECT con.body_my,con.game_id,con.channel_code,g.game_name,ch.name as channel_name,con.begin_time,con.end_time,con.state,"+
		"ou.incr_time,ou.recharge_time,ou.server_time "+
		"FROM contract con "+
		"LEFT JOIN game g ON con.game_id=g.game_id LEFT JOIN channel ch ON con.channel_code = ch.cp "+
		"LEFT JOIN game_outage ou ON con.game_id=ou.game_id "+
		"WHERE con.state=157 AND con.company_type=1 AND con.effective_state=1 %s ORDER BY con.game_id ", condition)
	var gameOutages []orm.Params

	_, err = orm.NewOrm().Raw(sql, dates...).Values(&gameOutages)
	if err != nil || len(gameOutages) == 0 {
		return
	}

	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//	keys := []string{}
	//	for _, contract := range gameOutages {
	//		gameId, _ := util.Interface2Int(contract["game_id"], false)
	//		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
	//		channelCode, _ := util.Interface2String(contract["cp"], false)
	//		days := (publishTime - time.Now().Unix()) / 3600 / 24
	//		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
	//		keys = append(keys, key)
	//	}
	//	oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"游戏名称", "我方主体", "合同状态", "渠道名称", "商务负责人", "签订/终止时间", "关闭新增时间", "关闭充值时间", "关闭服务器时间"}
	var contents [][]string
	for _, outage := range gameOutages {
		bodyId, _ := util.Interface2Int(outage["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		channelName, _ := util.Interface2String(outage["channel_name"], false)
		channelCode, _ := util.Interface2String(outage["channel_code"], false)
		gameName, _ := util.Interface2String(outage["game_name"], false)
		gameId, _ := util.Interface2Int(outage["game_id"], false)
		endTime, _ := util.Interface2Int(outage["end_time"], false)
		beginTime, _ := util.Interface2Int(outage["begin_time"], false)
		contractState := "即将停运"
		incrTime, _ := util.Interface2Int(outage["incr_time"], false)
		rechargeTime, _ := util.Interface2Int(outage["recharge_time"], false)
		serverTime, _ := util.Interface2Int(outage["server_time"], false)

		incrDays := TimeSubDays(time.Unix(incrTime, incrTime), time.Now())
		rechargeDays := TimeSubDays(time.Unix(rechargeTime, rechargeTime), time.Now())
		serverDays := TimeSubDays(time.Unix(serverTime, serverTime), time.Now())
		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}
		people := models.GetPeopleByChannelCode(channelCode)

		contents = append(contents, []string{
			gameName,
			bodyMy,
			contractState,
			channelName,
			people,
			fmt.Sprintf("%s~%s", time.Unix(beginTime, 0).Format("2006-01-02"), time.Unix(int64(endTime), 0).Format("2006-01-02")),
			time.Unix(incrTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(rechargeTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(serverTime, 0).Format("2006-01-02 15:04:05"),
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]即将停运!", gameName)
		if incrDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止新增!", incrDays)
		}
		if rechargeDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止充值!", rechargeDays)
		}
		if serverDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将关闭服务器!", serverDays)
		}
		info = info + fmt.Sprintf("与渠道[%s]的合同待处理!", channelName)

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			ChannelCode: channelName,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// 游戏停运cp合同待处理
func GameOutageCpContractOpe() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_OUTAGE_CP_OPERATE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(true, intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = " AND FROM_UNIXTIME(ou.server_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}

	sql := fmt.Sprintf("SELECT con.body_my,con.game_id,con.company_id,g.issue,g.game_name,com.name,con.begin_time,con.end_time,con.state,"+
		"ou.incr_time,ou.recharge_time,ou.server_time "+
		"FROM contract con "+
		"LEFT JOIN game g ON con.game_id=g.game_id LEFT JOIN company_type com ON g.issue = com.id "+
		"LEFT JOIN game_outage ou ON con.game_id=ou.game_id "+
		"WHERE con.state=157 AND con.company_type=0 AND con.effective_state=1 %s ORDER BY con.game_id ", condition)
	var gameOutages []orm.Params

	_, err = orm.NewOrm().Raw(sql, dates...).Values(&gameOutages)
	if err != nil || len(gameOutages) == 0 {
		return
	}

	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//	keys := []string{}
	//	for _, contract := range gameOutages {
	//		gameId, _ := util.Interface2Int(contract["game_id"], false)
	//		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
	//		channelCode, _ := util.Interface2String(contract["cp"], false)
	//		days := (publishTime - time.Now().Unix()) / 3600 / 24
	//		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
	//		keys = append(keys, key)
	//	}
	//	oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"游戏名称", "我方主体", "合同状态", "CP名称", "商务负责人", "签订/终止时间", "关闭新增时间", "关闭充值时间", "关闭服务器时间"}
	var contents [][]string
	for _, outage := range gameOutages {
		bodyId, _ := util.Interface2Int(outage["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		companyName, _ := util.Interface2String(outage["name"], false)
		gameName, _ := util.Interface2String(outage["game_name"], false)
		gameId, _ := util.Interface2Int(outage["game_id"], false)
		endTime, _ := util.Interface2Int(outage["end_time"], false)
		beginTime, _ := util.Interface2Int(outage["begin_time"], false)
		contractState := "即将停运"
		incrTime, _ := util.Interface2Int(outage["incr_time"], false)
		rechargeTime, _ := util.Interface2Int(outage["recharge_time"], false)
		serverTime, _ := util.Interface2Int(outage["server_time"], false)

		incrDays := TimeSubDays(time.Unix(incrTime, incrTime), time.Now())
		rechargeDays := TimeSubDays(time.Unix(rechargeTime, rechargeTime), time.Now())
		serverDays := TimeSubDays(time.Unix(serverTime, serverTime), time.Now())
		issue, _ := util.Interface2Int(outage["issue"], false)

		company, _ := models.GetDistributionCompanyById(int(issue))
		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}
		yunName := company.YunduanResPerName
		if yunName == "" {
			yunName = "无"
		}
		youName := company.YouliangResPerName
		if youName == "" {
			youName = "无"
		}

		contents = append(contents, []string{
			gameName,
			bodyMy,
			contractState,
			companyName,
			fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName),
			fmt.Sprintf("%s~%s", time.Unix(beginTime, 0).Format("2006-01-02"), time.Unix(int64(endTime), 0).Format("2006-01-02")),
			time.Unix(incrTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(rechargeTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(serverTime, 0).Format("2006-01-02 15:04:05"),
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]即将停运!", gameName)
		if incrDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止新增!", incrDays)
		}
		if rechargeDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将停止充值!", rechargeDays)
		}
		if serverDays >= 0 {
			info = info + fmt.Sprintf("还有[%v]天即将关闭服务器!", serverDays)
		}
		info = info + fmt.Sprintf("与CP[%s]的合同待处理!", companyName)

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			ChannelCode: companyName,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// 停运游戏渠道合同未下架
func GameOutageChannelContractDown() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_OUTAGE_CHANNEL_DOWN)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(false, intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = " AND FROM_UNIXTIME(ou.server_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}

	sql := fmt.Sprintf("SELECT con.body_my,con.game_id,con.channel_code,g.game_name,ch.name AS channel_name,con.begin_time,con.end_time,con.state,"+
		"ty.name AS state_name,ou.incr_time,ou.recharge_time,ou.server_time "+
		"FROM contract con "+
		"LEFT JOIN game g ON con.game_id=g.game_id LEFT JOIN channel ch ON con.channel_code=ch.cp "+
		"LEFT JOIN game_outage ou ON con.game_id=ou.game_id LEFT JOIN `types` ty ON con.state=ty.id "+
		"WHERE con.state!=155 AND con.company_type=1 AND con.effective_state=1 %s ORDER BY con.game_id ", condition)
	var gameOutages []orm.Params

	_, err = orm.NewOrm().Raw(sql, dates...).Values(&gameOutages)
	if err != nil || len(gameOutages) == 0 {
		return
	}

	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//	keys := []string{}
	//	for _, contract := range gameOutages {
	//		gameId, _ := util.Interface2Int(contract["game_id"], false)
	//		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
	//		channelCode, _ := util.Interface2String(contract["cp"], false)
	//		days := (publishTime - time.Now().Unix()) / 3600 / 24
	//		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
	//		keys = append(keys, key)
	//	}
	//	oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"游戏名称", "我方主体", "合同状态", "渠道名称", "商务负责人", "签订/终止时间", "关闭新增时间", "关闭充值时间", "关闭服务器时间"}
	var contents [][]string
	for _, outage := range gameOutages {
		bodyId, _ := util.Interface2Int(outage["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		channelName, _ := util.Interface2String(outage["channel_name"], false)
		channelCode, _ := util.Interface2String(outage["channel_code"], false)
		gameName, _ := util.Interface2String(outage["game_name"], false)
		gameId, _ := util.Interface2Int(outage["game_id"], false)
		endTime, _ := util.Interface2Int(outage["end_time"], false)
		beginTime, _ := util.Interface2Int(outage["begin_time"], false)
		contractState, _ := util.Interface2String(outage["state_name"], false)
		incrTime, _ := util.Interface2Int(outage["incr_time"], false)
		rechargeTime, _ := util.Interface2Int(outage["recharge_time"], false)
		serverTime, _ := util.Interface2Int(outage["server_time"], false)

		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}
		people := models.GetPeopleByChannelCode(channelCode)

		contents = append(contents, []string{
			gameName,
			bodyMy,
			contractState,
			channelName,
			people,
			fmt.Sprintf("%s~%s", time.Unix(beginTime, 0).Format("2006-01-02"), time.Unix(int64(endTime), 0).Format("2006-01-02")),
			time.Unix(incrTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(rechargeTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(serverTime, 0).Format("2006-01-02 15:04:05"),
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]已经关服，但是与渠道[%s]的合同还未下架!", gameName, channelName)

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			ChannelCode: channelName,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// 停运游戏CP合同未下架
func GameOutageCpContractDown() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_OUTAGE_CP_DOWN)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	dates := IntervalsToDates(false, intervals...)

	condition := ""
	if len(dates) > 0 {
		holder := strings.Repeat(",?", len(dates))
		condition = " AND FROM_UNIXTIME(ou.server_time,'%Y-%m-%d') in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		//return
	}

	sql := fmt.Sprintf("SELECT con.body_my,con.game_id,con.company_id,g.issue,g.game_name,com.name,con.begin_time,con.end_time,con.state,"+
		"ty.name AS state_name,ou.incr_time,ou.recharge_time,ou.server_time "+
		"FROM contract con "+
		"LEFT JOIN game g ON con.game_id=g.game_id LEFT JOIN company_type com ON g.issue = com.id "+
		"LEFT JOIN game_outage ou ON con.game_id=ou.game_id LEFT JOIN `types` ty ON con.state=ty.id "+
		"WHERE con.state!=155 AND con.company_type=0 AND con.effective_state=1 %s ORDER BY con.game_id ", condition)
	var gameOutages []orm.Params

	_, err = orm.NewOrm().Raw(sql, dates...).Values(&gameOutages)
	if err != nil || len(gameOutages) == 0 {
		return
	}

	//oldLog := map[string]bool{}
	//if rule.IsRepeat == 0 {
	//	keys := []string{}
	//	for _, contract := range gameOutages {
	//		gameId, _ := util.Interface2Int(contract["game_id"], false)
	//		publishTime, _ := util.Interface2Int(contract["publish_time"], false)
	//		channelCode, _ := util.Interface2String(contract["cp"], false)
	//		days := (publishTime - time.Now().Unix()) / 3600 / 24
	//		key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)
	//		keys = append(keys, key)
	//	}
	//	oldLog, _ = getRepeatWarningLog(rule.Type, keys)
	//}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"游戏名称", "我方主体", "合同状态", "CP名称", "商务负责人", "签订/终止时间", "关闭新增时间", "关闭充值时间", "关闭服务器时间"}
	var contents [][]string
	for _, outage := range gameOutages {
		bodyId, _ := util.Interface2Int(outage["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))
		companyName, _ := util.Interface2String(outage["name"], false)
		gameName, _ := util.Interface2String(outage["game_name"], false)
		gameId, _ := util.Interface2Int(outage["game_id"], false)
		endTime, _ := util.Interface2Int(outage["end_time"], false)
		beginTime, _ := util.Interface2Int(outage["begin_time"], false)
		contractState, _ := util.Interface2String(outage["state_name"], false)
		incrTime, _ := util.Interface2Int(outage["incr_time"], false)
		rechargeTime, _ := util.Interface2Int(outage["recharge_time"], false)
		serverTime, _ := util.Interface2Int(outage["server_time"], false)
		issue, _ := util.Interface2Int(outage["issue"], false)

		company, _ := models.GetDistributionCompanyById(int(issue))

		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, publishTime, days)

		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}
		yunName := company.YunduanResPerName
		if yunName == "" {
			yunName = "无"
		}
		youName := company.YouliangResPerName
		if youName == "" {
			youName = "无"
		}

		contents = append(contents, []string{
			gameName,
			bodyMy,
			contractState,
			companyName,
			fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName),
			fmt.Sprintf("%s~%s", time.Unix(beginTime, 0).Format("2006-01-02"), time.Unix(int64(endTime), 0).Format("2006-01-02")),
			time.Unix(incrTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(rechargeTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(serverTime, 0).Format("2006-01-02 15:04:05"),
		})

		// format the message and content
		info := fmt.Sprintf("请注意,[%s]已经关服，但是与CP[%s]的合同还未下架!", gameName, companyName)

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			GameId:      int(gameId),
			Grade:       rule.Grade,
			ChannelCode: companyName,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// 渠道主合同到期预警
func ChannelMainContractExpireWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_MAIN_CONTRACT_EXPIRE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}

	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	numbers := checker.IntervalsToDates(intervals...)
	condition := ""
	if len(numbers) > 0 {
		holder := strings.Repeat(",?", len(numbers))
		condition = "HAVING max_end_time in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		return
	}

	var contracts []orm.Params
	_, err = orm.NewOrm().Raw(fmt.Sprintf("SELECT a.company_id,a.begin_time,a.end_time AS max_end_time,a.body_my,b.channel_code "+
		"FROM main_contract a LEFT JOIN channel_company b ON a.company_id=b.company_id WHERE a.company_type=2 "+
		"%s", condition), numbers...).Values(&contracts)
	if err != nil || len(contracts) == 0 {
		return
	}

	//keys := []string{}
	//for _, contract := range contracts {
	//	endTime, _ := util.Interface2Int(contract["max_end_time"], false)
	//	channelCode, _ := util.Interface2String(contract["channel_code"], false)
	//	days := (endTime - time.Now().Unix()) / 3600 / 24
	//	key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)
	//	keys = append(keys, key)
	//}

	//oldLog, _ := getRepeatWarningLog(rule.Type, keys)
	userEmails := GetEmailsByUserIds(rule.UserIds)
	//userEmails := []string{"liuqilin@kuaifazs.com"}

	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "渠道名称", "公司名称", "商务负责人", "签订/终止时间"}
	var contents [][]string

	for _, contract := range contracts {
		beginTime, _ := util.Interface2String(contract["begin_time"], false)
		endTime, _ := util.Interface2String(contract["max_end_time"], false)
		channelCode, _ := util.Interface2String(contract["channel_code"], false)
		end, err := time.Parse("2006-01-02", endTime)
		if err != nil {
			continue
		}
		days := (end.Unix() - time.Now().Unix()) / 3600 / 24
		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)
		//
		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}

		//beginTime, _ := util.Interface2Int(contract["begin_time"], false)
		channelCompany, _ := models.GetChannelCompanyByCode(channelCode)
		bodyId, _ := util.Interface2Int(contract["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))

		companyName := ""
		if channelCompany != nil {
			companyName = channelCompany.CompanyName
		}

		info := fmt.Sprintf("请注意,[%s]-[%s]的主合同还有[%d]天到期",
			channelCompany.ChannelName,
			companyName,
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
			//game.GameName,
			fmt.Sprintf("<span style='color: red'>%s~%s</span>", beginTime, endTime),
		})

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			//GameId:      game.GameId,
			Grade:       rule.Grade,
			ChannelCode: channelCompany.ChannelCode,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}

// CP主合同到期预警
func CpMainContractExpireWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CP_MAIN_CONTRACT_EXPIRE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}

	rule := rules[0]

	intervals := StringToIntArray(rule.Intervals, ",")
	numbers := checker.IntervalsToDates(intervals...)
	condition := ""
	if len(numbers) > 0 {
		holder := strings.Repeat(",?", len(numbers))
		condition = "HAVING max_end_time in " + fmt.Sprintf("(%s)", holder[1:])
	} else {
		fmt.Errorf("请为%s预警设置提前预警时间", rule.Type)
		return
	}

	var contracts []orm.Params
	_, err = orm.NewOrm().Raw(fmt.Sprintf("SELECT a.company_id,a.begin_time,a.end_time AS max_end_time,a.body_my,"+
		"b.id FROM main_contract a LEFT JOIN distribution_company b ON a.company_id=b.company_id WHERE a.company_type=1 "+
		"%s", condition), numbers...).Values(&contracts)
	if err != nil || len(contracts) == 0 {
		return
	}

	//keys := []string{}
	//for _, contract := range contracts {
	//	endTime, _ := util.Interface2Int(contract["max_end_time"], false)
	//	channelCode, _ := util.Interface2String(contract["channel_code"], false)
	//	days := (endTime - time.Now().Unix()) / 3600 / 24
	//	key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)
	//	keys = append(keys, key)
	//}

	//oldLog, _ := getRepeatWarningLog(rule.Type, keys)
	userEmails := GetEmailsByUserIds(rule.UserIds)
	//userEmails := []string{"liuqilin@kuaifazs.com"}

	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	titles := []string{"我方主体", "公司名称", "商务负责人", "签订/终止时间"}
	var contents [][]string

	for _, contract := range contracts {
		beginTime, _ := util.Interface2String(contract["begin_time"], false)
		endTime, _ := util.Interface2String(contract["max_end_time"], false)
		distributionId, _ := util.Interface2Int(contract["id"], false)
		end, err := time.Parse("2006-01-02", endTime)
		if err != nil {
			continue
		}
		days := (end.Unix() - time.Now().Unix()) / 3600 / 24
		//key := fmt.Sprintf("%d|%s|%d|%d", gameId, channelCode, endTime, days)
		//
		//if rule.IsRepeat == 0 {
		//	if _, ok := oldLog[key]; ok {
		//		continue
		//	}
		//}

		//beginTime, _ := util.Interface2Int(contract["begin_time"], false)

		company, _ := models.GetDistributionCompanyById(int(distributionId))
		bodyId, _ := util.Interface2Int(contract["body_my"], false)
		bodyMy := ParseBodyMy(int(bodyId))

		companyName := ""
		if company != nil {
			companyName = company.CompanyName
		}

		info := fmt.Sprintf("请注意,CP[%s]的主合同还有[%d]天到期",
			companyName,
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
			companyName,
			fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName),
			//game.GameName,
			fmt.Sprintf("<span style='color: red'>%s~%s</span>", beginTime, endTime),
		})

		models.AddWarningLog(&models.WarningLog{
			WarningName: rule.Name,
			WarningType: rule.Type,
			Info:        info,
			CreateTime:  time.Now().Unix(),
			//GameId:      game.GameId,
			Grade:       rule.Grade,
			ChannelCode: company.CompanyName,
			//Keys:        key,
		})
	}

	if len(contents) > 0 {
		SendTableMail(subject, titles, contents, userEmails)
	}

	return
}
