package models

import (
	"github.com/astaxie/beego/orm"
	"errors"
	"time"
	"strconv"
	"fmt"
)

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

/*
 * 结算部统计:
 * 1)渠道寄出合同; 2)渠道归档合同; 3)渠道对账数; 4)CP寄出合同; 5) CP归档合同; 6)CP对账数
 */
type StatisticAccounting struct {
	Date                    string `json:"date"`
	ChannelSendContract     int    `json:"channel_send_contract"`     //寄出的渠道合同数量
	ChannelCompleteContract int    `json:"channel_complete_contract"` //归档的渠道合同数量
	ChannelVerify           int    `json:"channel_verify"`            //渠道对账数
	CpSendContract          int    `json:"cp_send_contract"`          //CP寄出合同
	CpCompleteContract      int    `json:"cp_complete_contract"`      //Cp归档合同数
	CpVerify                int    `json:"cp_verify"`                 //Cp对账数
}

// StatisticOfAccounting 返回结算部的统计详情
func StatisticOfAccounting(startTimestamp int64, endTimestamp int64, limit int, offset int) (statistics []*StatisticAccounting, total int64) {
	total = (endTimestamp - startTimestamp) / 86400
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	pageStart := end.AddDate(0, 0, -offset)
	pageEnd := pageStart.AddDate(0, 0, -limit)
	if pageEnd.Before(start) {
		pageEnd = start
	}
	pageStartUnix := pageStart.Unix()
	pageEndtUnix := pageEnd.Unix()

	if pageStartUnix == pageEndtUnix {
		pageEndtUnix = pageStartUnix - 86399
	}

	channelSend, _ := StatisticOfChannelSendContract(pageEndtUnix, pageStartUnix)
	channelComplete, _ := StatisticOfChannelCompleteContract(pageEndtUnix, pageStartUnix)
	channelVerify, _ := StatisticOfChannelVerify(pageEndtUnix, pageStartUnix)
	cpSend, _ := StatisticOfCpSendContract(pageEndtUnix, pageStartUnix)
	cpComplete, _ := StatisticOfCpCompleteContract(pageEndtUnix, pageStartUnix)
	cpVerify, _ := StatisticOfCpVerify(pageEndtUnix, pageStartUnix)

	statistics = []*StatisticAccounting{}

	fmt.Printf("start:%v---end:%v\n", pageStart, pageEnd)
	for i := pageStart; i.After(pageEnd.AddDate(0, 0, -1)); i = i.AddDate(0, 0, -1) {
		date := i.Format("2006-01-02")
		statistic := StatisticAccounting{}
		statistic.Date = date
		if channelSend[date] != nil {
			statistic.ChannelSendContract, _ = strconv.Atoi(channelSend[date].(string))
		}
		if channelComplete[date] != nil {
			statistic.ChannelCompleteContract, _ = strconv.Atoi(channelComplete[date].(string))
		}
		if channelVerify[date] != nil {
			statistic.ChannelVerify, _ = strconv.Atoi(channelVerify[date].(string))
		}
		if cpSend[date] != nil {
			statistic.CpSendContract, _ = strconv.Atoi(cpSend[date].(string))
		}
		if cpComplete[date] != nil {
			statistic.CpCompleteContract, _ = strconv.Atoi(cpComplete[date].(string))
		}
		if cpVerify[date] != nil {
			statistic.CpVerify, _ = strconv.Atoi(cpVerify[date].(string))
		}
		statistics = append(statistics, &statistic)
	}

	return
}

// StatisticOfChannelSendContract 1)渠道寄出合同: 根据起始时间和结束时间，统计该时间段的渠道寄出的合同数量
func StatisticOfChannelSendContract(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	start := time.Unix(startTimestamp, 0).Format("2006-01-02")
	end := time.Unix(endTimestamp, 0).Format("2006-01-02")
	sql := "SELECT `time` as `key`, SUM(`send`) AS `value` " +
		"FROM `contract_statistics` " +
		"WHERE `time` >= ? AND `time` <= ? AND `type` = 2 " +
		"GROUP BY `time` " +
		"ORDER BY `time` desc"
	result, err = getKeyValues(sql, start, end)
	return
}

// StatisticOfChannelCompleteContract  2)渠道归档合同: 根据起始时间和结束时间，统计该时间段每日合同归档数量
func StatisticOfChannelCompleteContract(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	start := time.Unix(startTimestamp, 0).Format("2006-01-02")
	end := time.Unix(endTimestamp, 0).Format("2006-01-02")
	sql := "SELECT `time` as `key`, SUM(`complete`) AS `value` " +
		"FROM `contract_statistics` " +
		"WHERE `time` >= ? AND `time` <= ? AND `type` = 2 " +
		"GROUP BY `time` " +
		"ORDER BY `time` desc"
	result, err = getKeyValues(sql, start, end)
	return
}

// StatisticOfCpVerify 3)渠道对账数：统计指定时间范围内每天添加的渠道对账条数
func StatisticOfChannelVerify(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(created_time, '%Y-%m-%d') AS `key`, COUNT(*) AS `value` " +
		"FROM verify_channel " +
		" WHERE created_time >= ? AND created_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"

	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticOfCpSendContract 4)CP寄出合同: 根据起始时间和结束时间，统计该时间段的渠道寄出的合同数量
func StatisticOfCpSendContract(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	start := time.Unix(startTimestamp, 0).Format("2006-01-02")
	end := time.Unix(endTimestamp, 0).Format("2006-01-02")
	sql := "SELECT `time` as `key`, SUM(`send`) AS `value` " +
		"FROM `contract_statistics` " +
		"WHERE `time` >= ? AND `time` <= ? AND `type` = 1 " +
		"GROUP BY `time` " +
		"ORDER BY `time` desc"
	result, err = getKeyValues(sql, start, end)
	return
}

// StatisticOfChannelCompleteContract  5)CP归档合同: 根据起始时间和结束时间，统计该时间段每日合同归档数量
func StatisticOfCpCompleteContract(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	start := time.Unix(startTimestamp, 0).Format("2006-01-02")
	end := time.Unix(endTimestamp, 0).Format("2006-01-02")
	sql := "SELECT `time` as `key`, SUM(`complete`) AS `value` " +
		"FROM `contract_statistics` " +
		"WHERE `time` >= ? AND `time` <= ? AND `type` = 1 " +
		"GROUP BY `time` " +
		"ORDER BY `time` desc"
	result, err = getKeyValues(sql, start, end)
	return
}

// StatisticOfCpVerify 6)CP对账数：统计指定时间范围内每天添加的CP对账条数，正确则返回KeyValue数组
func StatisticOfCpVerify(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(created_time, '%Y-%m-%d') AS `key`, COUNT(*) AS `value` " +
		"FROM verify_cp " +
		"WHERE created_time >= ? AND created_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"

	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

/*
 * 财务部统计:
 * 1)渠道回款; 2)CP付款
 */

type StatisticFinancial struct {
	Date        string  `json:"date"`
	ChannelPaid float64 `json:"channel_paid"` //渠道回款
	CpPaid      float64 `json:"cp_paid"`      //CP回款
}

// StatisticOfFinancial 按分页返回财务部的数据统计
func StatisticOfFinancial(startTimestamp int64, endTimestamp int64, limit int, offset int) (statistics []*StatisticFinancial, total int64) {
	total = (endTimestamp - startTimestamp) / 86400
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	pageStart := end.AddDate(0, 0, -offset)
	pageEnd := pageStart.AddDate(0, 0, -limit)
	if pageEnd.Before(start) {
		pageEnd = start
	}
	pageStartUnix := pageStart.Unix()
	pageEndtUnix := pageEnd.Unix()

	if pageStartUnix == pageEndtUnix {
		pageEndtUnix = pageStartUnix - 86399
	}

	channelPaid, _ := StatisticOfChannelPaid(pageEndtUnix, pageStartUnix)
	cpPaid, _ := StatisticOfCpPaid(pageEndtUnix, pageStartUnix)
	statistics = []*StatisticFinancial{}
	for i := pageStart; i.After(pageEnd.AddDate(0, 0, -1)); i = i.AddDate(0, 0, -1) {
		date := i.Format("2006-01-02")
		statistic := StatisticFinancial{}
		statistic.Date = date
		if channelPaid[date] != nil {
			statistic.ChannelPaid, _ = strconv.ParseFloat(channelPaid[date].(string), 64)
		}
		if cpPaid[date] != nil {
			statistic.CpPaid, _ = strconv.ParseFloat(cpPaid[date].(string), 64)
		}

		statistics = append(statistics, &statistic)
	}
	return
}

// StatisticOfChannelPaid 1)渠道回款: 统计起始时间戳到结束时间戳之间的渠道回款情况, 返回KeyValue数组
func StatisticOfChannelPaid(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(remit_time, '%Y-%m-%d') as `key`, SUM(amount) AS `value` " +
		"FROM `remit_down_account` " +
		"WHERE remit_time >= ? AND remit_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticOfCpPaid 2)CP付款: 统计起始时间戳到结束时间戳之间的CP结算情况
func StatisticOfCpPaid(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(settle_time, '%Y-%m-%d') AS `key`, SUM(amount) AS `value` " +
		"FROM `settle_down_account` " +
		"WHERE settle_time >= ? AND settle_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

/*
 * 渠道商务:
 * 1)渠道接入; 2)渠道回款(同财务部渠道回款); 3)渠道未回款(暂时无法统计); 4)渠道商信息
 */

type StatisticChannelTrade struct {
	Date           string  `json:"date"`
	ChannelAccess  int     `json:"channel_access"`  // 渠道接入数量
	ChannelPaid    float64 `json:"channel_paid"`    // 渠道回款金额
	ChannelCompany int     `json:"channel_company"` // 新添加渠道商数量
}

// StatisticOfChannelTrade 按分页统计渠道商务的统计数据
func StatisticOfChannelTrade(startTimestamp int64, endTimestamp int64, limit int, offset int) (statistics []*StatisticChannelTrade, total int64) {
	total = (endTimestamp - startTimestamp) / 86400
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	pageStart := end.AddDate(0, 0, -offset)
	pageEnd := pageStart.AddDate(0, 0, -limit)
	if pageEnd.Before(start) {
		pageEnd = start
	}
	pageStartUnix := pageStart.Unix()
	pageEndtUnix := pageEnd.Unix()

	if pageStartUnix == pageEndtUnix {
		pageEndtUnix = pageStartUnix - 86399
	}

	channelAccess, _ := StatisticOfChannelAccess(pageEndtUnix, pageStartUnix)
	channelPaid, _ := StatisticOfChannelPaid(pageEndtUnix, pageStartUnix)
	channelCompany, _ := StatisticOfChannelCompany(pageEndtUnix, pageStartUnix)
	statistics = []*StatisticChannelTrade{}
	for i := pageStart; i.After(pageEnd.AddDate(0, 0, -1)); i = i.AddDate(0, 0, -1) {
		date := i.Format("2006-01-02")
		statistic := StatisticChannelTrade{}
		statistic.Date = date
		if channelAccess[date] != nil {
			statistic.ChannelAccess, _ = strconv.Atoi(channelAccess[date].(string))
		}
		if channelPaid[date] != nil {
			statistic.ChannelPaid, _ = strconv.ParseFloat(channelPaid[date].(string), 64)
		}
		if channelCompany[date] != nil {
			statistic.ChannelCompany, _ = strconv.Atoi(channelCompany[date].(string))
		}

		statistics = append(statistics, &statistic)
	}
	return
}

// StatisticOfChannelAccess 1)渠道接入: 根据起始时间和结束时间，计算该时间范围内每天的渠道接入数量
func StatisticOfChannelAccess(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	result, err = getNewRowsPerDay("channel_access", startTimestamp, endTimestamp)
	return
}

// TODO:按天统计是否存在问题，并且对账也是针对月份的
// StatisticOfChannelNotPay 渠道商务-渠道未回款，统计某个时间段内每天的未回款金额
func StatisticOfChannelNotPay(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT c.body_my, c.remit_company_id, d.name, IFNULL(a.amount, 0) - IFNULL(b.amount, 0) AS total_amount " +
		"FROM remit_pre_account AS c " +
		"LEFT JOIN ( " +
		"SELECT body_my,remit_company_id,SUM(amount_payable) AS amount " +
		"FROM `verify_channel` WHERE `status` >20 " +
		"GROUP BY body_my,remit_company_id) AS a ON c.remit_company_id = a.remit_company_id AND c.body_my = a.body_my " +
		"LEFT JOIN ( " +
		"SELECT body_my,remit_company_id,SUM(amount) AS amount " +
		"FROM `remit_down_account` " +
		"GROUP BY body_my,remit_company_id) AS b ON c.remit_company_id = b.remit_company_id AND a.body_my = b.body_my "

	var data []orm.Params
	_, err = orm.NewOrm().Raw(sql, startTimestamp, endTimestamp).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}
	result = map[interface{}]interface{}{}
	for _, per := range data {
		result[per["settle_time"]] = per["amount"]
	}
	return
}

// StatisticOfChannelCompany 4)渠道商信息：根据起始时间和结束时间，计算该间范围内的每天新增渠道商数量
func StatisticOfChannelCompany(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	result, err = getNewRowsPerDay("channel_company", startTimestamp, endTimestamp)
	return
}

/*
 * CP商务:
 * 1)提测游戏; 2)接入游戏; 3)CP付款(同财务部CP付款); 4)发行商信息
 */

type StatisticCpTrade struct {
	Date                string  `json:"date"`
	GameEvaluate        int     `json:"game_evaluate"`        // 提测游戏数量
	GameAccess          int     `json:"game_access"`          // 接入游戏数量
	CpPaid              float64 `json:"cp_paid"`              // CP付款金额
	DistributionCompany int     `json:"distribution_company"` // 新添加发行商数量
}

// StatisticOfCTrade 按分页统计CP商务部门的数据统计
func StatisticOfCpTrade(startTimestamp int64, endTimestamp int64, limit int, offset int) (statistics []*StatisticCpTrade, total int64) {
	fmt.Println(limit, offset)
	total = (endTimestamp - startTimestamp) / 86400
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	pageStart := end.AddDate(0, 0, -offset)
	pageEnd := pageStart.AddDate(0, 0, -limit)
	if pageEnd.Before(start) {
		pageEnd = start
	}
	pageStartUnix := pageStart.Unix()
	pageEndtUnix := pageEnd.Unix()

	if pageStartUnix == pageEndtUnix {
		pageEndtUnix = pageStartUnix - 86399
	}

	gameEvaluate, _ := StatisticOfGameEvaluate(pageEndtUnix, pageStartUnix)
	gameAccess, _ := StatisticOfGameAccess(pageEndtUnix, pageStartUnix)
	cpPaid, _ := StatisticOfCpPaid(pageEndtUnix, pageStartUnix)
	distributionCompany, _ := StatisticOfDistributionCompany(pageEndtUnix, pageStartUnix)
	statistics = []*StatisticCpTrade{}
	for i := pageStart; i.After(pageEnd.AddDate(0, 0, -1)); i = i.AddDate(0, 0, -1) {
		date := i.Format("2006-01-02")
		statistic := StatisticCpTrade{}
		statistic.Date = date
		if gameEvaluate[date] != nil {
			statistic.GameEvaluate, _ = strconv.Atoi(gameEvaluate[date].(string))
		}
		if gameAccess[date] != nil {
			statistic.GameAccess, _ = strconv.Atoi(gameAccess[date].(string))
		}
		if cpPaid[date] != nil {
			statistic.CpPaid, _ = strconv.ParseFloat(cpPaid[date].(string), 64)
		}
		if distributionCompany[date] != nil {
			statistic.DistributionCompany, _ = strconv.Atoi(distributionCompany[date].(string))
		}
		statistics = append(statistics, &statistic)
	}
	return
}

// StatisticOfCpPaid 1)提测游戏: 根据开始时间和结束时间，计算该段时间内每天CP提测的游戏数量
func StatisticOfGameEvaluate(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(create_time, '%Y-%m-%d') AS `key`, count(*) AS `value` " +
		"FROM `game` " +
		"WHERE game_id = 0 AND create_time >= ? AND create_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticOfCpPaid 2)接入游戏: 根据开始时间和结束时间，计算该段时间内每天CP接入的游戏数量
func StatisticOfGameAccess(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(import_time, '%Y-%m-%d') AS `key`, count(*) AS `value` " +
		"FROM `game` " +
		"WHERE game_id > 0 AND import_time >= ? AND import_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticOfDistributionCompany 4)发行商信息: 根据开始时间和结束时间，计算指定时间范围内的每天新增发行商数量
func StatisticOfDistributionCompany(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	result, err = getNewRowsPerDay("distribution_company", startTimestamp, endTimestamp)
	return
}

/*
 * 运营部:
 * 1)通过游戏; 2)未通过游戏; 3)运营准备
 */

type StatisticOperation struct {
	Date          string `json:"date"`
	AgreedGame    int    `json:"agreed_game"`     // 评测通过的游戏总数
	NotAgreedGame int    `json:"not_agreed_game"` // 评测未通过的游戏总数
	Prepare       int    `json:"papare"`          // 新添加的运营负责人总数
}

// StatisticOfFinancial 按分页返回财务部的数据统计
func StatisticOfOperation(startTimestamp int64, endTimestamp int64, limit int, offset int) (statistics []*StatisticOperation, total int64) {
	total = (endTimestamp - startTimestamp) / 86400
	start := time.Unix(startTimestamp, 0)
	end := time.Unix(endTimestamp, 0)

	pageStart := end.AddDate(0, 0, -offset)
	pageEnd := pageStart.AddDate(0, 0, -limit)
	if pageEnd.Before(start) {
		pageEnd = start
	}
	pageStartUnix := pageStart.Unix()
	pageEndtUnix := pageEnd.Unix()

	if pageStartUnix == pageEndtUnix {
		pageEndtUnix = pageStartUnix - 86399
	}

	agreedGame, _ := StatisticOfAgreedGame(pageEndtUnix, pageStartUnix)
	notAgreedGame, _ := StatisticOfNotAgreedGame(pageEndtUnix, pageStartUnix)
	prepare, _ := StatisticOfPrepare(pageEndtUnix, pageStartUnix)
	statistics = []*StatisticOperation{}
	for i := pageStart; i.After(pageEnd.AddDate(0, 0, -1)); i = i.AddDate(0, 0, -1) {
		date := i.Format("2006-01-02")
		statistic := StatisticOperation{}
		statistic.Date = date
		if agreedGame[date] != nil {
			statistic.AgreedGame, _ = strconv.Atoi(agreedGame[date].(string))
		}
		if notAgreedGame[date] != nil {
			statistic.NotAgreedGame, _ = strconv.Atoi(notAgreedGame[date].(string))
		}
		if prepare[date] != nil {
			statistic.Prepare, _ = strconv.Atoi(prepare[date].(string))
		}

		statistics = append(statistics, &statistic)
	}
	return
}

// StatisticOfAgreedGame 1)通过的游戏: 根据起始时间和终止时间，计算该段时间内评测通过的游戏总数
func StatisticOfAgreedGame(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(result_time, '%Y-%m-%d') AS `key`, count(*) AS `value` " +
		"FROM `game` " +
		"WHERE result_time >= ? AND result_time <= ? AND `advise` LIKE '%:1%' " +
		"GROUP BY `key` " +
		"ORDER BY `key`"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticONotAgreedGame 2)未通过的游戏: 根据起始时间和终止时间，计算该段时间内评测未通过的游戏总数
func StatisticOfNotAgreedGame(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(result_time, '%Y-%m-%d') AS `key`, count(*) AS `value` " +
		"FROM `game` " +
		"WHERE result_time >= ? AND result_time <= ? AND `advise` NOT LIKE '%:1%' " +
		"GROUP BY `key` " +
		"ORDER BY `key`"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// StatisticOfPrepare 3)新添加的运营负责人总数: 根据起始时间和终止时间，计算该段时间内运营负责人的总数
func StatisticOfPrepare(startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	sql := "SELECT FROM_UNIXTIME(operator_time, '%Y-%m-%d') AS `key`, sum(LENGTH(operator_person) - LENGTH(REPLACE(operator_person,',',''))+ 1) AS `value` " +
		"FROM `game_plan` " +
		"WHERE operator_time >= ? AND operator_time <= ? AND operator_person IS NOT NULL AND LENGTH(operator_person)>1 " +
		"GROUP BY `key` " +
		"ORDER BY `key`"
	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// GetNewRowsPerDay 获取指定时间范围内某个表每天的新增记录
func getNewRowsPerDay(tableName string, startTimestamp int64, endTimestamp int64) (result map[interface{}]interface{}, err error) {
	result = map[interface{}]interface{}{}
	if tableName == "" {
		err = errors.New("table name is null")
		return
	}
	sql := "SELECT FROM_UNIXTIME(create_time, '%Y-%m-%d') AS `key`, COUNT(*) AS `value` " +
		"FROM " + tableName +
		" WHERE create_time >= ? AND create_time <= ? " +
		"GROUP BY `key` " +
		"ORDER BY `key` desc"

	result, err = getKeyValues(sql, startTimestamp, endTimestamp)
	return
}

// GetKeyValues 根据传入的sql，返回key/value数组
func getKeyValues(sql string, start interface{}, end interface{}) (result map[interface{}]interface{}, err error) {
	var data []orm.Params
	_, err = orm.NewOrm().Raw(sql, start, end).Values(&data)
	if err != nil || len(data) == 0 {
		return
	}
	result = map[interface{}]interface{}{}
	for _, per := range data {
		//result = append(result, &KeyValue{Key: per["key"], Value: per["value"]})
		result[per["key"]] = per["value"]
	}
	return
}
