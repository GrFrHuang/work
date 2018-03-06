package tool

import (
	"bytes"
	"github.com/bysir-zl/bygo/util"
	"text/template"
	"github.com/bysir-zl/bygo/log"
	"kuaifa.com/kuaifa/work-together/models"
)

const (
	sub_appUpdate = "欢迎使用 Work Together"
	con_appUpdate = "{{.name}}({{.nickname}}) 你好! 你的初始密码为：{{.pwd}}．请尽快修改. 登陆Work Together: http://work.kuaifazs.com"
	sub_resetPwd  = "[Work Together] 重置密码成功"
	con_resetPwd  = "{{.name}}({{.nickname}}) 你好! 你的密码为：{{.pwd}}．请尽快修改. "
	sub_workflow  = "[Work Together] 工作流"
	con_workflow  = "{{.nickname}} 你好! 你的[Work Together] 工作流为：{{.name}}. 还未处理! 请尽快处理. 登陆Work Together: http://work.kuaifazs.com"
)

var (
	mail *util.Mail
	tpl  *template.Template = template.New("singleton")
)

func SendInitPwdEmail(data map[string]string, emailAddr string) error {
	subject := tplString(sub_appUpdate, data)
	content := tplString(con_appUpdate, data)
	return mail.Send(subject, content, []string{emailAddr})
}
func SendResetPwdEmail(name, nickname, pwd string, emailAddr string) error {
	data := map[string]string{
		"name":     name,
		"nickname": nickname,
		"pwd":      pwd,
	}

	subject := tplString(sub_resetPwd, data)
	content := tplString(con_resetPwd, data)
	return mail.Send(subject, content, []string{emailAddr})
}

func SendWorkflowEmail(nickname,taskName,emailAddr string) error {
	data := map[string]string{
		"nickname": nickname,
		"name": taskName,
	}
	subject := tplString(sub_workflow, data)
	content := tplString(con_workflow, data)
	return mail.Send(subject, content, []string{emailAddr})
}

func tplString(temp string, data map[string]string) string {
	t, err := tpl.Parse(temp)
	if err != nil {
		log.Error("tpl mail", err)
		return ""
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		log.Error("tpl mail", err)
	}
	return buf.String()
}

func SendMessageEmail(subject, body string, to []string) error {
	err := mail.Send(subject, body, to)
	return err
}

func SendMessageEmailByDp(subject, body string, deId int) {
	users, _ := models.GetUsersByDevMent(deId)
	toUsers := []string{}
	for _, user := range users {
		toUsers = append(toUsers, user.Email)

	}
	SendMessageEmail(subject, body, toUsers)
}

func init() {
	mail = util.NewMail("kuaifazhushou@163.com", "wushuang", "smtp.163.com:25")
}
