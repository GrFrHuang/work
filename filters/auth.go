package filters

import (
	"github.com/astaxie/beego/context"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strings"
)

// 除了登录的api(/v1/session), 其他api都需要登录
var AuthLogin = func(ctx *context.Context) {
	if strings.Contains(ctx.Request.RequestURI, "/session") || strings.Contains(ctx.Request.RequestURI, "/face") ||
		strings.Contains(ctx.Request.RequestURI, "/sign") || strings.Contains(ctx.Request.RequestURI, "/userfaceverify") {
		return
	}
	if ctx.Request.Method == "OPTIONS" {
		ctx.Input.SetData("uid", 0)
		return
	}
	err, uid := tool.VerifyToken(ctx)
	if err != nil {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                               //允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST,DELETE, GET, PUT, OPTIONS") //允许post访问
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")     //header的类型
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
		ctx.ResponseWriter.ResponseWriter.WriteHeader(bean.CODE_Unauthorized)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Input.SetData("uid", uid)
}
