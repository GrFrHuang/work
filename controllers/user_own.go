package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// 用户自己的信息
type UserOwnController struct {
	BaseController
}

func (c *UserOwnController) URLMapping() {
	c.Mapping("Put", c.Put)
}

// ChangePwd
// @Title 自己修改自己密码
// @Description 自己修改自己密码
// @Param	body		body 	models.bean.InputPwd	true		"body for pwd"
// @Success 200
// @Failure 403 body is empty
// @router /pwd [put]
func (c *UserOwnController) Put() {
	var v bean.InputPwd
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err)
		return
	}
	pwd := v.NewPwd
	if l := len(pwd); l < 6 || l > 20 {
		c.RespJSON(bean.CODE_Params_Err, "密码长度[6,20]")
		return
	}

	token := tool.GetRequestToken(c.Ctx)
	token = token
	u, err := models.GetUserInfoByToken(token);
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	err = tool.ChangeOwnPwd(u.Name, v.OldPwd, v.NewPwd, token)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// @Title 自己修改自己的信息
// @Description 自己修改自己的信息
// @Param	body		body 	models.User	true		""
// @Success 200
// @Failure 403 body is empty
// @router /info [put]
func (c *UserOwnController) ModifyInfo() {
	v := models.User{Id: c.Uid()}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("Unmarshal err: %v", err)
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if err := models.UpdateUserSelfDataById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// @Title 获取自己的信息+权限
// @Description 获取自己信息
// @Success 200
// @Failure 403 body is empty
// @router / [get]
func (c *UserOwnController) Get() {
	uid := c.Uid()
	u, err := models.GetUserById(uid)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	d, err := models.GetDepartmentById(u.DepartmentId)
	if err == nil {
		u.DepartmentName = d.Name
	}

	rsMap := models.GetAllRoleMap()
	models.AddRoleInfo(u,rsMap)

	menu,isAdmin:=models.GetCanVisitMenu(uid)
	result:=map[string]interface{}{
		"info":u,
		"menu":menu,
		"isAdmin":isAdmin,
	}
	c.RespJSONData(result)
}
