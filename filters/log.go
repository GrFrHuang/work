package filters

import (
	"github.com/astaxie/beego/context"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/tool"
	"sort"
	"strings"
)

var logText = map[string]string{
	"get /logs":   "-",
	"get /asset/": "-",

	"post /session":   "登录",
	"delete /session": "注销",

	"get /department":    "查看部门",
	"delete /department": "删除部门",
	"post /department":   "修改部门",
	"put /department":    "添加部门",

	"get /permission": "查看权限列表",

	"get /user":         "查看用户",
	"delete /user":      "删除用户",
	"put /user":         "修改用户",
	"post /user":        "添加用户",
	"get /user/devment": "-",

	"get /role":    "查看角色",
	"delete /role": "删除角色",
	"put /role":    "修改角色",
	"post /role":   "添加角色",

	"get /game": "查看游戏",

	"get /contract":    "查看合同",
	"put /contract":    "修改合同",
	"delete /contract": "删除合同",

	"get /order":         "查看流水",
	"get /order/game":    "-",
	"get /order/channel": "-",

	"get /remitaccount":         "查看回款单",
	"get /remitaccount/pre":     "查看未回款单",
	"get /remitaccount/channel": "-",

	"get /settleaccount":      "查看结算单",
	"post /settleaccount":     "添加结算账单",
	"get /settleaccount/pre":  "查看未结算单",
	"get /settleaccount/game": "-",

	"get /channelverify":                "查看渠道对账",
	"post /channelverify":               "添加渠道对账",
	"put /channelverify":                "修改渠道对账",
	"get /channelverify/remitcompanies": "-",
	"get /channelverify/novdate":        "-",

	"post /cpverify": "添加CP对账单",

	"get /types": "查看类型分组",
	// todo add more
}
var sortedLog tool.StringLenSort = []string{}

var Logger = func(ctx *context.Context) {
	if ctx.Request.Method == "options" {
		return
	}
	url := ctx.Request.Method + " " + strings.Replace(ctx.Request.RequestURI, "/v1", "", -1)
	url = strings.ToLower(url)

	logName := url
	for _, v := range sortedLog {
		if strings.Index(url, v) == 0 {
			name := logText[v]
			logName = name
			break
		}
	}
	if logName == "-" {
		return
	}

	uid := 0
	if u := ctx.Input.GetData("uid"); u != nil {
		uid = u.(int)
	}
	if uid == 0 {
		return
	}
	if !strings.Contains(url, "face") && !strings.Contains(url, "sign") {
		input := "url:" + url + ",body:" + string(ctx.Input.RequestBody)
		models.SaveOperation(uid, logName, ctx.ResponseWriter.Status, input)
	}
}

func init() {
	for k := range logText {
		sortedLog = append(sortedLog, k)
	}
	// 根据url长度排序, 优先匹配长的
	sort.Sort(sortedLog)
	sort.Reverse(sortedLog)
}
