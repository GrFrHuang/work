package warning

import (
	"github.com/bysir-zl/bygo/util"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models"
	"encoding/json"
)

var Mail *util.Mail

func init() {
	Mail = util.NewMail("kuaifazhushou@163.com", "wushuang", "smtp.163.com:25")
}

func SendTableMail(subject string, titles []string, contents [][]string, recipients []string) {
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
	htmlcontent += `<div><br>具体请前往：<a href="http://work.kuaifazs.com/#/login" target="_blank">Work Together</a>查看详情!</div>`
	err := Mail.Send(subject, htmlcontent, recipients)
	if err != nil {
		beego.Warning(err.Error())
	}
}

func GetEmailsByUserIdsBak(ids string) (emails []string) {
	emails = []string{}
	id_list := strings.Split(ids, ",")
	for _, id := range id_list {
		b, _ := strconv.ParseInt(id, 10, 32)
		user, _ := models.GetUserById(int(b))
		if user == nil {
			continue
		}
		emails = append(emails, user.Email)
	}
	return
}

func GetEmailsByUserIds(ids string) (emails []string) {
	if ids == "" {
		return
	}

	type user struct {
		Type  int //表示个人或者整个部门，1：部门，2：个人
		Value interface{}
		name  string
	}
	emails = []string{}
	var users []user
	fmt.Println(ids)
	err := json.Unmarshal([]byte(ids), &users)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range users {
		if v.Type == 1 {
			id := int(v.Value.(float64))
			dev_users, err := models.GetUsersByDevMent(id)
			if err == nil {
				for _, i := range dev_users {
					emails = append(emails, i.Email)
				}
			}
		} else if v.Type == 2 {
			emails = append(emails, v.Value.(string))
		} else {
			fmt.Println("error type")
		}
	}
	return
}
