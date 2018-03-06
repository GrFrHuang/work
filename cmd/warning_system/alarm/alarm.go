package alarm

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"time"
)

const (
	verify_account_context   = "[%s]渠道%s %s未对账, 总流水 [%.2f]" // 360,皇室战争:1000.0,startTime-now,total
	remit_account_context    = "[%s]主体%s未回款, 总金额 [%.2f]"    // 360,startTime-now,total
	contract_timeout_context = "[%s]渠道[%s]游戏的合同%s"
	contract_sign_context    = "[%s]渠道[%s]游戏%s, 但还未签订合同"
	order                    = ""
)

func SaveRemitAlarmLog(remitCompanyId int, startTime int, total float64) (err error) {
	company := models.Company{Id: remitCompanyId}
	err = orm.NewOrm().Read(&company)
	if err != nil {
		return
	}
	t := time.Unix(int64(startTime), 0)
	nowTimestamp := time.Now().Unix()
	dayInt := (nowTimestamp - int64(startTime)) / 3600 / 24
	day := strconv.FormatInt(dayInt, 10)

	l := &models.AlarmLog{
		Context: fmt.Sprintf(remit_account_context,
			company.Name,
			fmt.Sprintf("[%s]至今(共[%s]天)",t.Format("2006-01-02"),day),
			total,
		),
		CreatedTime: int(nowTimestamp),
		IsHide:      2,
		Type:        "remit_account_context",
	}
	_, err = models.AddAlarmLog(l)
	return
}

func SaveVerifyAlarmLog(channelCode string, startTime int, games string, total float64) (err error) {
	channel := models.Channel{Cp: channelCode}
	err = orm.NewOrm().Read(&channel, "Cp")
	if err != nil {
		return
	}
	t := time.Unix(int64(startTime), 0)
	nowTimestamp := time.Now().Unix()
	dayInt := (nowTimestamp - int64(startTime)) / 3600 / 24
	day := strconv.FormatInt(dayInt, 10)

	l := &models.AlarmLog{
		Context: fmt.Sprintf(verify_account_context,
			channel.Name,
			games,
			fmt.Sprintf("[%s]至今(共[%s]天)",t.Format("2006-01-02"),day),
			total,
		),
		CreatedTime: int(nowTimestamp),
		IsHide:      2,
		Type:        "verify_account_context",
	}
	_, err = models.AddAlarmLog(l)
	return
}

func SaveContractTimeout(channelName string, gameName string, days int64) (err error) {

	nowTimestamp := time.Now().Unix()
	day := ""
	if days < 0 {
		day = fmt.Sprintf("已经过期[%d]天", -days)
	} else {
		day = fmt.Sprintf("将于[%d]天后过期", days)
	}
	l := &models.AlarmLog{
		Context: fmt.Sprintf(contract_timeout_context,
			channelName,
			gameName,
			day,
		),
		CreatedTime: int(nowTimestamp),
		IsHide:      2,
		Type:        "contract_timeout_context",
	}
	_, err = models.AddAlarmLog(l)
	return
}

func SaveContractNotSign(channelName string, gameName string, days int64) (err error) {

	nowTimestamp := time.Now().Unix()
	day := ""
	if days < 0 {
		day = fmt.Sprintf("已经发行[%d]天", -days)
	} else {
		day = fmt.Sprintf("将于[%d]天后发行", days)
	}
	l := &models.AlarmLog{
		Context: fmt.Sprintf(contract_sign_context,
			channelName,
			gameName,
			day,
		),
		CreatedTime: int(nowTimestamp),
		IsHide:      2,
		Type:        "contract_sign_context",
	}
	_, err = models.AddAlarmLog(l)
	return
}
