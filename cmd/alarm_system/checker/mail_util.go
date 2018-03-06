package checker

import (
	"fmt"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"github.com/astaxie/beego"
)

var RECIVER_LIST = []string{"zhouqiangwei@kuaifazs.com", "yangfei@kuaifazs.com", "chenfang@kuaifazs.com",
                            "leihanjian@kuaifazs.com", "fengxia@kuaifazs.com", "lannana@kuaifazs.com", "wuchengshuang@kuaifazs.com"}

var RECIVER_LIST_DEBUG = []string{"zhangliang@kuaifazs.com"}

var Mail *util.Mail

func init() {
	Mail = util.NewMail("feichen_noreply@163.com", "feichen123", "smtp.163.com:25")
}

func SendRemitAccount(remit NotRemit) {
	subject := fmt.Sprintf("【渠道未回款提醒】%s", remit.RemitCompanyName)
	content := fmt.Sprintf("未回款金额：%f<br/>未回款时间：%s<br/><br/>需了解详情请前往work together系统，http://work.kuaifazs.com", remit.Amount, time.Unix(remit.StartTime, 0).Format("2006-01-02"))
	recipients := RECIVER_LIST
	Mail.Send(subject, content, recipients)
}

func SendVerifyAccount(verify NotVerify) {
	subject := fmt.Sprintf("【渠道未回款提醒】%s（%s）", verify.CompanyName, verify.ChannelName)
	content := fmt.Sprintf("未对账金额：%f<br/>未对账时间：%s<br/><br/>需了解详情请前往work together系统，http://work.kuaifazs.com",
		verify.Amount,
		time.Unix(int64(verify.StartTime), 0).Format("2006-01-02"))
	recipients := RECIVER_LIST
	Mail.Send(subject, content, recipients)
}

func SendTableMail(subject string, titles []string, contents [][]string) {
	var htmlcontent string
	htmlcontent = `<table style="width:100%;" cellpadding="2" cellspacing="0" align="center" border="1" bordercolor="#000000">` + `<tbody><tr style="background-color:#337FE5;">`
	for _, title := range titles {
		htmlcontent += fmt.Sprintf(`<td>%s</td>`, title)
	}
	htmlcontent += `</tr>`
	for _, content := range contents {
		htmlcontent += `<tr>`
		for _, c := range content {
			htmlcontent += fmt.Sprintf(`<td>%s</td>`, c)
		}
		htmlcontent += `</tr>`
	}
	htmlcontent += `</tbody></table>`
	recipients := RECIVER_LIST
	err := Mail.Send(subject, htmlcontent, recipients)
	if err != nil {
		beego.Warning(err.Error())
	}

}

func SendContractTimeOut(n models.Contract) {
	body := ""
	sign_time := time.Unix(n.SigningTime, 0).Format("2006-01-02 15:04:05")
	end_time := time.Unix(n.EndTime, 0).Format("2006-01-02 15:04:05")
	if n.BodyMy == 1 {
		body = "云端"
	} else if n.BodyMy == 2 {
		body = "有量"
	}
	content := fmt.Sprintf("合同签订时间：%s<br/>合同终止时间：%s<br/>我方主体：%s<br/>合作游戏：%s<br/><br/>需了解详情请前往work together系统，http://work.kuaifazs.com",
		sign_time, end_time, body, n.Game.GameName)
	subject := fmt.Sprintf("【渠道合同到期提醒】%s(%s)", n.ChannelCompanyName, n.Channel.Name)

	recipients := RECIVER_LIST
	Mail.Send(subject, content, recipients)
}

func SendContractNotSign(companyName, gameName string) {
	subject := fmt.Sprintf("【渠道合同未签订提醒】%s", companyName)
	content := fmt.Sprintf("合作游戏：%s<br/><br/>需了解详情请前往work together系统，http://work.kuaifazs.com", gameName)
	recipients := RECIVER_LIST
	err := Mail.Send(subject, content, recipients)
	if err != nil {
		log.Error("mail", err)
	}
}
