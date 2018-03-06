package checker

import (
	"github.com/bysir-zl/bygo/util"
	"strings"
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"strconv"
)

// SendGameUpdateWarning 当游戏需要更新时,添加更新后根据传入的gameId
// 发送预警和更新预警日志
func SendGameUpdateWarning(gameId int) (err error) {
	o := orm.NewOrm()
	game := models.Game{GameId: gameId}
	err = o.Read(&game, "GameId")
	if err != nil {
		return
	}

	gameUpdate := models.GameUpdate{GameId: gameId}
	err = o.Read(&gameUpdate, "GameId")
	if err != nil {
		return
	}

	rules, err := models.GetWarningByType(models.WARNING_GAME_UPDATE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	updateType := ""
	if gameUpdate.UpdateType == 1 {
		updateType = "强更"
	} else if gameUpdate.UpdateType == 2 {
		updateType = "热更"
	} else {
		updateType = "错误类型"
	}

	userEmails := GetEmailsByUserIds(rule.UserIds)
	titles := []string{"游戏名称", "更新类型", "更新时间"}
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	updateTimeFormat := time.Unix(gameUpdate.GameUpdateTime, 0).Format("2006-01-02")
	contents := [][]string{{game.GameName, updateType, updateTimeFormat}}
	info := fmt.Sprintf("[%s]将与[%s]进行[%s]", game.GameName, updateTimeFormat, updateType)

	models.AddWarningLog(&models.WarningLog{
		WarningName: rule.Name,
		WarningType: rule.Type,
		Info:        info,
		CreateTime:  time.Now().Unix(),
		Keys:        strconv.Itoa(game.GameId),
		Grade:       rule.Grade})

	SendTableMail(subject, titles, contents, userEmails)
	return
}

// getRepeatWarningLog 查询预警日志中是否已经存在当前报警, 如果已经存在
// 则设置该key为true
func getRepeatWarningLog(warningType string, keys []string) (oldKeys map[string]bool, err error) {
	cond := []interface{}{}
	oldKeys = map[string]bool{}
	cond = append(cond, warningType)
	for _, key := range keys {
		cond = append(cond, key)
	}

	condition := ""
	if len(keys) > 0 {
		holder := strings.Repeat(",?", len(keys))
		condition = fmt.Sprintf(" AND `keys` in(%s) ", holder[1:])
	}

	logs := []orm.Params{}
	sql_log := "SELECT `keys` FROM `warning_log` WHERE `warning_type`=? " + condition
	_, err = orm.NewOrm().Raw(sql_log, cond...).Values(&logs)
	if err != nil {
		return
	}

	for _, log := range logs {
		key, _ := util.Interface2String(log["keys"], false)
		oldKeys[key] = true
	}

	return
}

// SendNewGameAccessWarning 当有新游戏接入时调用该方法可以直接发送gameId相关的
// 预警消息和更新预警日志
func SendNewGameAccessWarning(gameId int) (err error) {
	o := orm.NewOrm()
	game := models.Game{GameId: gameId}
	err = o.Read(&game, "GameId")
	if err != nil {
		return
	}

	rules, err := models.GetWarningByType(models.WARNING_GAME_ACCESS)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}

	rule := rules[0]
	user_emails := GetEmailsByUserIds(rule.UserIds)
	titles := []string{"游戏名称", "首发时间", "接入时间"}
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))

	importTimeFormat := time.Unix(game.ImportTime, 0).Format("2006-01-02")
	releaseTimeFormat := time.Unix(game.PublishTime, 0).Format("2006-01-02")
	contents := [][]string{{game.GameName, releaseTimeFormat, importTimeFormat}}

	info := fmt.Sprintf("[%s]于[%s]完成游戏接入，将于[%s]首发", game.GameName, importTimeFormat, releaseTimeFormat)
	models.AddWarningLog(&models.WarningLog{
		WarningName: rule.Name,
		WarningType: rule.Type,
		Info:        info,
		CreateTime:  time.Now().Unix(),
		Keys:        strconv.Itoa(game.GameId),
		GameId:	     game.GameId,
		Grade:       rule.Grade})

	SendTableMail(subject, titles, contents, user_emails)
	return
}

// SendGameReleaseUpdateWarning 当游戏首发时间被修改时根据传入的gameId, oldTime, newTime
// 发送预警和更新预警日志
func SendGameReleaseUpdateWarning(gameId int, oldTime int64, newTime int64) (err error) {
	o := orm.NewOrm()
	game := models.Game{GameId: gameId}
	err = o.Read(&game, "GameId")
	if err != nil {
		return
	}

	rules, err := models.GetWarningByType(models.WARNING_GAME_RELEASE_UPDATE)
	if err != nil || len(rules) == 0 || rules[0].State == 0 {
		return
	}
	rule := rules[0]

	user_emails := GetEmailsByUserIds(rule.UserIds)
	titles := []string{"游戏名称", "旧首发时间", "新首发时间"}
	subject := fmt.Sprintf("%s(%s预警)", rule.Type, models.ParseWarningGrade(rule.Grade))
	oldTimeFormat := time.Unix(oldTime, 0).Format("2006-01-02")
	newTimeFormat := time.Unix(newTime, 0).Format("2006-01-02")
	contents := [][]string{{game.GameName, oldTimeFormat, newTimeFormat}}
	info := fmt.Sprintf("[%s]将首发时间[%s]调整为[%s]", game.GameName, oldTimeFormat, newTimeFormat)

	models.AddWarningLog(&models.WarningLog{
		WarningName: rule.Name,
		WarningType: rule.Type,
		Info:        info,
		CreateTime:  time.Now().Unix(),
		Keys:        strconv.Itoa(game.GameId),
		Grade:       rule.Grade})

	SendTableMail(subject, titles, contents, user_emails)
	return
}

// TimestampsToDates 根据所给天数间隔获取未来的日期
// 字符串数组
func IntervalsToDates(days ...int) (dates []interface{}) {
	now := time.Now()
	dates = []interface{}{}
	for _, v := range days {
		dates = append(dates, now.AddDate(0, 0, v).Format("2006-01-02"))
	}
	return
}

// StringToIntArray 将以sep分割的int字符串转转换成为
// int数组
func StringToIntArray(intString string, sep string) (arrInt []int) {
	arrString := strings.Split(intString, sep)
	arrInt = []int{}
	for _, v := range arrString {
		i, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		arrInt = append(arrInt, i)
	}
	return
}

// ParseBodyMy 将我方主体Id编号转换为对应的主体名称
// 1 -> "云端"， 2 -> "有量"
func ParseBodyMy(bodyId int) (bodyMy string) {
	if bodyId == 1 {
		bodyMy = "云端"
	} else if bodyId == 2 {
		bodyMy = "有量"
	} else {
		bodyMy = "错误的主体类型"
	}
	return
}