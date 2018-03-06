package checker

import (
	"kuaifa.com/kuaifa/work-together/models"
	"github.com/astaxie/beego/orm"
	"time"
	"github.com/bysir-zl/bygo/util"
	"strings"
	"fmt"
)

func NewGameAccessWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_ACCESS)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		games := []orm.Params{}
		_, err = orm.NewOrm().Raw("SELECT `game_id`, `game_name`, `import_time`, `publish_time` " +
				"FROM game WHERE `publish_time`>? AND `publish_time`<? ORDER BY `import_time` DESC", time.Now().Unix(), time.Now().Unix()).Values(&games)

		//_, err = orm.NewOrm().Raw("SELECT `game_id`, `game_name`, `import_time`, `publish_time` " +
		//		"FROM game WHERE `publish_time`>? AND `publish_time`<? ORDER BY `import_time` DESC", 1431619200, 1437408000).Values(&games)
		//fmt.Println(games)

		if err != nil || len(games) == 0 {
			return
		}

		// 找到已经发送的game
		old_games := map[int]bool{}
		if rule.IsRepeat == 0 {
			game_ids := []int64{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_GAME_ACCESS)
			for _, game := range games {
				game_id, _ := util.Interface2Int(game["game_id"], false)
				game_ids = append(game_ids, game_id)
				cond = append(cond, game_id)
			}

			condition := ""

			if len(game_ids) > 0 {
				holder := strings.Repeat(",?", len(games))
				condition = fmt.Sprintf(" AND game_id in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT `game_id` FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2Int(log["game_id"], false)
				old_games[int(id)] = true
			}
		}

		for _, game := range games {
			game_id, _ := util.Interface2Int(game["game_id"], false)
			if rule.IsRepeat == 0 {
				if _, ok := old_games[int(game_id)]; ok {
					continue
				}
			}

			// format the message and content
			gameName, _ := util.Interface2String(game["game_name"], false)
			importTime, _ := util.Interface2Int(game["import_time"], false)
			releaseTime, _ := util.Interface2Int(game["publish_time"], false)

			content := fmt.Sprintf("[%s]于[%s]完成游戏接入，将于[%s]首发",
				gameName,
				time.Unix(int64(importTime), 0).Format("2006-01-02"),
				time.Unix(int64(releaseTime), 0).Format("2006-01-02"))
			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, GetEmailsByUserIds(rule.Emails))
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{
				WarningName: rule.Name,
				WarningType: rule.Type,
				Info:content,
				CreateTime:time.Now().Unix(),
				GameId:int(game_id),
				Grade: rule.Grade})
		}
	}
	return
}

func ReleaseTimeUpdateWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_RELEASE_UPDATE)
	if err != nil || len(rules) == 0 {
		return
	}

	//for _, rule := range rules {
	//
	//}

	return
}

func GameUpdateWarning() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_GAME_UPDATE)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}
		o := orm.NewOrm()
		v := &models.GameUpdate{}
		games := []models.GameUpdate{}
		_, err = o.QueryTable(v).Filter("game_update_time__gt", time.Now().Unix()).Filter("game_update_time__lt", time.Now().Unix()).All(&games)
		//_, err = o.QueryTable(v).Filter("game_update_time__gt", 1488988700).Filter("game_update_time__lt", 1490025601).All(&games)
		if err != nil || len(games) == 0 {
			return
		}

		old_games := map[int]bool{}
		if rule.IsRepeat == 0 {
			game_ids := []int{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_GAME_UPDATE)
			for _, game := range games {
				game_ids = append(game_ids, game.GameId)
				cond = append(cond, game.GameId)
			}

			condition := ""
			if len(game_ids) > 0 {
				holder := strings.Repeat(",?", len(games))
				condition = fmt.Sprintf(" AND game_id in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT `game_id` FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2Int(log["game_id"], false)
				old_games[int(id)] = true
			}
		}

		for _, game := range games {
			if rule.IsRepeat == 0 {
				if _, ok := old_games[game.GameId]; ok {
					continue
				}
			}

			// format the message and content
			update_type := "未知"
			if game.UpdateType == 1 {
				update_type = "强更"
			} else if game.UpdateType == 2 {
				update_type = "热更"
			} else {

			}

			tmp_game, _ := models.GetGameByGameId(game.Id)

			content := fmt.Sprintf("[%s]将于[%s]进行[%s]",
				tmp_game.GameName,
				time.Unix(int64(game.GameUpdateTime), 0).Format("2006-01-02"),
				update_type)

			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, GetEmailsByUserIds(rule.Emails))
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{WarningName: rule.Name, WarningType: rule.Type, Info:content, CreateTime:time.Now().Unix(), GameId:game.GameId, Grade: rule.Grade})
		}
	}
	return
}

func ChannelContractExpireWarningBak() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_CONTRACT_EXPIRE)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		contracts := []orm.Params{}
		_, err = orm.NewOrm().Raw("SELECT game_id, company_id, max_end_time FROM (SELECT game_id, company_id, MAX(end_time) AS max_end_time FROM contract WHERE company_type=1 GROUP BY game_id, company_id ) AS t WHERE max_end_time >? AND max_end_time < ?",
			time.Now().Unix(), time.Now().Unix()).Values(&contracts)

		//_, err = orm.NewOrm().Raw("SELECT game_id, company_id, max_end_time FROM (SELECT game_id, company_id, MAX(end_time) AS max_end_time FROM contract WHERE company_type=1 GROUP BY game_id, company_id ) AS t WHERE max_end_time >? AND max_end_time < ?",
		//	1489593600, 1490198400).Values(&contracts)
		if err != nil || len(contracts) == 0 {
			return
		}

		old_contract := map[string]bool{}

		// 如果已经报警这找出已经报警的游戏和渠道
		if rule.IsRepeat == 0 {
			game_channels := []string{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_CHANNEL_CONTRACT_EXPIRE)
			for _, con := range contracts {
				game_id, _ := util.Interface2String(con["game_id"], false)
				channel_id, _ := util.Interface2Int(con["company_id"], false)
				channel, _ := models.GetChannelByChannelId(int(channel_id))
				game_channel := fmt.Sprintf("%s%s", game_id, channel.Cp)
				fmt.Printf("channel:%v, game_id:%s, channle_code:%s, game_channel:%s", channel, game_id, channel.Cp, game_channel)
				game_channels = append(game_channels, game_channel)
				cond = append(cond, game_channel)
			}

			condition := ""
			if len(game_channels) > 0 {
				holder := strings.Repeat(",?", len(game_channels))
				condition = fmt.Sprintf(" AND concat(game_id, channel_code) in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT concat(`game_id`, `channel_code`) as game_channel FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2String(log["game_channel"], false)
				old_contract[id] = true
			}
		}

		fmt.Printf("old_channel:%v\n", old_contract)

		for _, contract := range contracts {
			game_id, _ := util.Interface2Int(contract["game_id"], false)
			end_time, _ := util.Interface2Int(contract["max_end_time"], false)
			channel_id, _ := util.Interface2Int(contract["company_id"], false)
			channel, _ := models.GetChannelByChannelId(int(channel_id))
			game_channel := fmt.Sprintf("%d%s", game_id, channel.Cp)

			if rule.IsRepeat == 0 {
				if _, ok := old_contract[game_channel]; ok {
					continue
				}
			}

			game, _ := models.GetGameByGameId(int(game_id))

			days := (end_time - time.Now().Unix()) / 3600 / 24

			// format the message and content
			content := fmt.Sprintf("请注意,[%s]-[%s]的合同还有[%d]天到期",
				game.GameName,
				channel.Name,
				days)

			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, GetEmailsByUserIds(rule.Emails))
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{WarningName: rule.Name, WarningType: rule.Type, Info:content, CreateTime:time.Now().Unix(), GameId:game.GameId, Grade: rule.Grade, ChannelCode:channel.Cp})
		}
	}
	return
}

func ChannelContractSign() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_CONTRACT_SIGN)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		contracts := []orm.Params{}
		_, err = orm.NewOrm().Raw("" +
				"SELECT g.`game_id`, g.`game_name`, g.`game_name`, ch.`cp`, g.`publish_time`, ch.`name` " +
				"FROM contract AS c LEFT JOIN game AS g  ON g.`game_id` = c.`game_id` LEFT JOIN channel AS ch ON c.company_id=ch.channel_id " +
				"WHERE c.`company_type`=1 AND c.`state`=0 AND g.`publish_time` >? AND g.`publish_time`<? ",
			time.Now().Unix(), time.Now().Unix()).Values(&contracts)

		//_, err = orm.NewOrm().Raw("" +
		//		"SELECT g.`game_id`, g.`game_name`, g.`game_name`, ch.`cp`, g.`publish_time`, ch.`name` " +
		//		"FROM contract AS c LEFT JOIN game AS g  ON g.`game_id` = c.`game_id` LEFT JOIN channel AS ch ON c.company_id=ch.channel_id " +
		//		"WHERE c.`company_type`=1 AND c.`state`=0 AND g.`publish_time` >? AND g.`publish_time`<? ",
		//	1490140800, 1490918400).Values(&contracts)

		if err != nil || len(contracts) == 0 {
			return
		}

		old_contract := map[string]bool{}

		// 如果已经报警这找出已经报警的游戏和渠道
		if rule.IsRepeat == 0 {
			game_channels := []string{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_CHANNEL_CONTRACT_SIGN)
			for _, con := range contracts {
				game_id, _ := util.Interface2String(con["game_id"], false)
				channel_code, _ := util.Interface2String(con["cp"], false)
				game_channel := fmt.Sprintf("%s%s", game_id, channel_code)
				game_channels = append(game_channels, game_channel)
				cond = append(cond, game_channel)
			}

			condition := ""
			if len(game_channels) > 0 {
				holder := strings.Repeat(",?", len(game_channels))
				condition = fmt.Sprintf(" AND concat(game_id, channel_code) in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT concat(`game_id`, `channel_code`) as game_channel FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2String(log["game_channel"], false)
				old_contract[id] = true
			}
		}

		user_emails := GetEmailsByUserIds(rule.Emails)
		for _, contract := range contracts {
			game_id, _ := util.Interface2Int(contract["game_id"], false)
			game_name, _ := util.Interface2String(contract["game_name"], false)
			channel_code, _ := util.Interface2String(contract["cp"], false)
			channel_name, _ := util.Interface2String(contract["name"], false)
			end_time, _ := util.Interface2Int(contract["publish_time"], false)
			game_channel := fmt.Sprintf("%d%s", game_id, channel_code)

			if rule.IsRepeat == 0 {
				if _, ok := old_contract[game_channel]; ok {
					continue
				}
			}

			days := (end_time - time.Now().Unix()) / 3600 / 24

			// format the message and content
			content := fmt.Sprintf("请注意,[%s]还有[%d]天首发，与渠道[%s]的合同还未签订",
				game_name,
				days,
				channel_name)

			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, user_emails)
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{
				WarningName: rule.Name,
				WarningType: rule.Type,
				Info:content,
				CreateTime:time.Now().Unix(),
				GameId:int(game_id),
				Grade: rule.Grade,
				ChannelCode:channel_code})
		}
	}
	return
}

func ChannelOrderVerify() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_VERIFY)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		data := []orm.Params{}
		_, err = orm.NewOrm().Raw("" +
				"SELECT o.game_id , g.game_name, o.channel_code, con.body_my, ch.`name` AS channel_name,SUM(o.amount) AS total,MIN(o.`date`) AS min_date " +
				"FROM `order_pre_verify_channel` AS o LEFT JOIN channel AS ch ON o.`channel_code` = ch.`cp` " +
				"LEFT JOIN contract AS con ON con.company_id = ch.channel_id AND  con.company_type = 1 AND  o.game_id = con.game_id " +
				"LEFT JOIN game AS g ON o.game_id = g.game_id " +
				"WHERE o.verify_id = 0 AND o.`date` < ? " +
				"GROUP BY o.game_id, o.channel_code, con.body_my " +
				"HAVING total>?",
			time.Now().AddDate(0, 0, 1).Format("2006-01"), rule.Amount).Values(&data)

		if err != nil || len(data) == 0 {
			return
		}

		// 如果已经报警这找出已经报警的游戏和渠道
		// old_data := map[string]bool{}
		/*if rule.IsRepeat == 0 {
			game_channel_dates := []string{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_CHANNEL_ORDER_VERIFY)
			for _, con := range data {
				game_id, _ := util.Interface2String(con["game_id"], false)
				channel_code, _ := util.Interface2String(con["channel_code"], false)
				date, _ := util.Interface2String(con["min_date"], false)
				game_channel_date := fmt.Sprintf("%s%s%s", game_id, channel_code, date)
				game_channel_dates = append(game_channel_dates, game_channel_date)
				cond = append(cond, game_channel_date)
			}

			condition := ""
			if len(game_channel_dates) > 0 {
				holder := strings.Repeat(",?", len(game_channel_dates))
				condition = fmt.Sprintf(" AND concat(`game_id`, `channel_code`, `date`) in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT concat(`game_id`, `channel_code`, `date`) as game_channel_date FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2String(log["game_channel_date"], false)
				old_data[id] = true
			}
		}*/

		user_emails := GetEmailsByUserIds(rule.Emails)
		for _, per := range data {
			game_id, _ := util.Interface2Int(per["game_id"], false)
			channel_code, _ := util.Interface2String(per["channel_code"], false)
			date, _ := util.Interface2String(per["min_date"], false)
			total, _ := util.Interface2Float(per["total"], false)

			channel_name, _ := util.Interface2String(per["channel_name"], false)
			game_name, _ := util.Interface2String(per["game_name"], false)

			//game_channel_date := fmt.Sprintf("%d%s%s", game_id, channel_code, date)
			//
			//if rule.IsRepeat == 0 {
			//	if _, ok := old_data[game_channel_date]; ok {
			//		continue
			//	}
			//}

			// format the message and content
			content := fmt.Sprintf("请注意,[%s]下的[%s]流水达到[%f]，最后一次未对账时间[%s]",
				channel_name,
				game_name,
				total,
				date)

			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, user_emails)
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{
				WarningName: rule.Name,
				WarningType: rule.Type,
				Info: content,
				CreateTime: time.Now().Unix(),
				GameId: int(game_id),
				Grade: rule.Grade,
				Amount: total,
				Date: date,
				ChannelCode: channel_code})
		}
	}
	return
}

func ChannelOrderPay() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_PAY)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		data := []orm.Params{}
		_, err = orm.NewOrm().Raw("" +
				"SELECT remit_company_id, MIN(date) AS start_time, SUM(amount_payable - amount_remit) AS amount " +
				"FROM verify_channel " +
				"WHERE amount_payable != amount_remit AND STATUS = 30 " +
				"GROUP BY remit_company_id;",
			time.Now().AddDate(0, 0, 1).Format("2006-01"), rule.Amount).Values(&data)

		if err != nil || len(data) == 0 {
			return
		}

		// 如果已经报警这找出已经报警的游戏和渠道
		// old_data := map[string]bool{}
		/*if rule.IsRepeat == 0 {
			game_channel_dates := []string{}
			cond := []interface{}{}
			cond = append(cond, models.WARNING_CHANNEL_ORDER_PAY)
			for _, con := range data {
				game_id, _ := util.Interface2String(con["game_id"], false)
				channel_code, _ := util.Interface2String(con["channel_code"], false)
				date, _ := util.Interface2String(con["min_date"], false)
				game_channel_date := fmt.Sprintf("%s%s%s", game_id, channel_code, date)
				game_channel_dates = append(game_channel_dates, game_channel_date)
				cond = append(cond, game_channel_date)
			}

			condition := ""
			if len(game_channel_dates) > 0 {
				holder := strings.Repeat(",?", len(game_channel_dates))
				condition = fmt.Sprintf(" AND concat(`game_id`, `channel_code`, `date`) in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT concat(`game_id`, `channel_code`, `date`) as game_channel_date FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				id, _ := util.Interface2String(log["game_channel_date"], false)
				old_data[id] = true
			}
		}*/

		user_emails := GetEmailsByUserIds(rule.Emails)
		for _, per := range data {
			game_id, _ := util.Interface2Int(per["game_id"], false)
			channel_code, _ := util.Interface2String(per["channel_code"], false)
			date, _ := util.Interface2String(per["min_date"], false)
			total, _ := util.Interface2Float(per["total"], false)

			channel_name, _ := util.Interface2String(per["channel_name"], false)
			game_name, _ := util.Interface2String(per["channel_name"], false)

			//game_channel_date := fmt.Sprintf("%d%s%s", game_id, channel_code, date)
			//
			//if rule.IsRepeat == 0 {
			//	if _, ok := old_data[game_channel_date]; ok {
			//		continue
			//	}
			//}

			// format the message and content
			content := fmt.Sprintf("请注意,[%s]下的[%s]流水达到[%f]，最后一次未回款时间[%s]",
				channel_name,
				game_name,
				total,
				date)

			subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)

			err := Mail.Send(subject, content, user_emails)
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{
				WarningName: rule.Name,
				WarningType: rule.Type,
				Info: content,
				CreateTime: time.Now().Unix(),
				GameId: int(game_id),
				Grade: rule.Grade,
				Amount: total,
				Date: date,
				ChannelCode: channel_code})
		}
	}
	return
}

func ChannelOrderThreshold() (err error) {
	rules, err := models.GetWarningByType(models.WARNING_CHANNEL_ORDER_THRESHOLD)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		if rule.State == 0 {
			continue
		}

		data := []orm.Params{}
		//_, err = orm.NewOrm().Raw("" +
		//		"SELECT o.game_id, o.cp AS channel_code, g.`game_name`, c.`name` AS channel_name, o.total " +
		//		"FROM( SELECT cp, game_id, SUM(amount) AS total FROM `order` " +
		//		"GROUP BY game_id, cp) AS o LEFT JOIN `game` AS g ON o.`game_id` = g.`game_id` LEFT JOIN channel AS c ON o.cp = c.`cp` " +
		//		"WHERE o.total > ? ", rule.Amount).Values(&data)

		_, err = orm.NewOrm().Raw("" +
				"SELECT o.game_id, o.cp AS channel_code, g.`game_name`, c.`name` AS channel_name, o.total " +
				"FROM( SELECT cp, game_id, SUM(amount) AS total FROM `order` " +
				"GROUP BY game_id, cp) AS o LEFT JOIN `game` AS g ON o.`game_id` = g.`game_id` LEFT JOIN channel AS c ON o.cp = c.`cp` " +
				"WHERE o.total > ? limit 5", 100).Values(&data)

		if err != nil || len(data) == 0 {
			return
		}

		 // 如果已经报警这找出已经报警的游戏和渠道
		 old_data := map[string]bool{}
		if rule.IsRepeat == 0 {
			cond := []interface{}{}
			cond = append(cond, models.WARNING_CHANNEL_ORDER_THRESHOLD)
			for _, con := range data {
				game_id, _ := util.Interface2String(con["game_id"], false)
				channel_code, _ := util.Interface2String(con["channel_code"], false)
				total, _ := util.Interface2Float(con["total"], false)
				key := fmt.Sprintf("%s%s%.2f", game_id, channel_code, total)
				cond = append(cond, key)
			}

			fmt.Println(cond)
			condition := ""
			if len(cond) - 1 > 0 {
				holder := strings.Repeat(",?", len(cond) - 1)
				condition = fmt.Sprintf(" AND concat(`game_id`, `channel_code`, `amount`) in(%s) ", holder[1:])
			}

			logs := []orm.Params{}
			sql_log := "SELECT `game_id`, `channel_code`, `amount` FROM `warning_log` WHERE `warning_type`=? " + condition
			_, err = orm.NewOrm().Raw(sql_log, cond).Values(&logs)
			for _, log := range logs {
				game_id, _ := util.Interface2String(log["game_id"], false)
				channel_code, _ := util.Interface2String(log["channel_code"], false)
				total, _ := util.Interface2Float(log["amount"], false)
				key := fmt.Sprintf("%s%s%.2f", game_id, channel_code, total)
				old_data[key] = true
			}
		}

		user_emails := GetEmailsByUserIds(rule.Emails)
		subject := fmt.Sprintf("【%s】预警-【%s-%s】", models.ParseWarningGrade(rule.Grade), rule.Type, rule.Name)
		titles := []string{"渠道名称", "游戏名称", "流水", "阈值"}
		contents := [][]string{}
		for _, per := range data {
			game_id, _ := util.Interface2Int(per["game_id"], false)
			channel_code, _ := util.Interface2String(per["channel_code"], false)
			total, _ := util.Interface2Float(per["total"], false)
			channel_name, _ := util.Interface2String(per["channel_name"], false)
			game_name, _ := util.Interface2String(per["game_name"], false)

			key := fmt.Sprintf("%d%s%0.2f", game_id, channel_code, total)
			if rule.IsRepeat == 0 {
				if _, ok := old_data[key]; ok {
					continue
				}
			}

			// format the message and content
			info := fmt.Sprintf("请注意,[%s]下的[%s]流水达到[%.2f]，已经超过[%.2f]阈值",
				channel_name,
				game_name,
				total,
				rule.Amount)
			contents = append(contents, []string{channel_name, game_name, fmt.Sprintf("%.2f", total), fmt.Sprintf("%.2f", rule.Amount)})

			//err := Mail.Send(subject, content, user_emails)
			if err != nil {
				continue
			}
			models.AddWarningLog(&models.WarningLog{
				WarningName: rule.Name,
				WarningType: rule.Type,
				Info: info,
				CreateTime: time.Now().Unix(),
				GameId: int(game_id),
				Grade: rule.Grade,
				Amount: total,
				ChannelCode: channel_code})
		}
		if len(contents) > 0 {
			SendTableMail(subject, titles, contents, user_emails)
		}
	}
	return
}