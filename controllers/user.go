package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
	"strings"
)

// 用户编辑
type UserController struct {
	BaseController
}

// URLMapping ...
func (c *UserController) URLMapping() {

}

// Post
// @Title 添加用户
// @Description create User
// @Param	body		body 	models.swagger.CreateUser	true		"body for User content"
// @ Param	x-token		header 	string	true		"x-token 测试时使用"
// @Success 200 {object} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	var ids []int
	if err := json.Unmarshal([]byte(v.RoleIds), &ids); err != nil {
		c.RespJSON(bean.CODE_Params_Err, "RoleIds is not json string")
		return
	}

	if err := c.Validate(v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	// 暂时用户名为邮箱号
	v.Name = strings.Split(v.Email, "@")[0]
	// 随机密码
	v.Pwd = tool.RandString(6)
	v.Disabled = 1

	resUser, err := tool.RegisterUser(v.Name, v.Pwd, v.Email, v.Phone, v.Nickname)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	v.Name = resUser.Name
	v.Email = resUser.Email
	if _, err := models.AddUser(&v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	go func(v models.User) {
		err = tool.SendInitPwdEmail(map[string]string{
			"pwd":      v.Pwd,
			"nickname": v.Nickname,
			"name":     v.Name,
		}, v.Email)
		if err != nil {
			log.Error("EMAIL", "发送邮件失败: "+err.Error())
		}
	}(v)

	v.Pwd = ""
	c.RespJSONData(v)
}

// GetOne
// @Title 获取user
// @Description get User by id
// @Param	id	path 	string	true		"The key for static block"
// @Success 200 {object} models.swagger.OutPutUser
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.swagger.OutPutUser
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
	var us []models.User
	f, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	f.Where["!name"] = []interface{}{""}
	f.Order = append(f.Order, "desc")
	f.Sortby = append(f.Sortby, "disabled")
	total, err := tool.GetAllByFilterWithTotal(new(models.User), &us, f)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	models.GroupAddRoleInfo(us)
	models.GroupAddDepartmentInfo(us)

	c.RespJSONDataWithTotal(us, total)
}

// Put ...
// @Title 更新用户
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.swagger.UpdateUserData	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.User{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("Unmarshal err: %v", err)
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.DepartmentId == 246 {
		v.Nickname = v.Nickname + "(已离职)"
	}
	if err := models.UpdateUserDataById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// @Title 禁止/恢复登陆
// @Description 禁止/恢复登陆
// @Param	disabled	query 	int	true		"1:恢复登陆 2:禁止登陆"
// @Success 200 {string} success!
// @router /disable/:id [put]
func (c *UserController) Disable() {
	disabled, _ := c.GetInt("disabled", 0)
	if disabled == 0 {
		c.RespJSON(bean.CODE_Params_Err, "disable is 0")
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	orm.NewOrm().Update(&models.User{Id: id, Disabled: disabled}, "Disabled")

	if disabled == 2 {
		models.EmptyUserToken(id)
	}

	c.RespJSONData("OK")
	return
}

// 修改密码
// @Title 管理员修改用户密码
// @Description 管理员修改用户密码
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.swagger.AdminChangePwd	true		"body for User content"
// @Success 200 {string} "OK"
// @router /resetpwd/:id [put]
func (c *UserController) ResetPwd() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	pwd := tool.RandString(6)
	u, err := models.GetUserById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	err = tool.ChangeUserPwd(u.Name, pwd)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	go func() {
		err := tool.SendResetPwdEmail(u.Name, u.Nickname, pwd, u.Email)
		if err != nil {
			log.Error("send email", err)
		}
	}()

	c.RespJSONData("OK")
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if id == 1 {
		// super super manager can not delete
		c.RespJSON(bean.CODE_Forbidden, "super manager can not delete")
		return
	}
	v, err := models.GetUserById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	err = tool.DeleteUser(v.Name)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	if err := models.DeleteUser(id); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// 这里暂时不设置权限
// @Title 获取用户组
// @Description 根据部门获取用户组
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} models.User
// @Failure 403 id is empty
// @router /devment/:id [get]
func (c *UserController) GetByDevMent() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	us, err := models.GetUsersByDevMent(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(us)
}

// @Title 获取当前用户所属部门
// @Description 获取当前用户所属部门
// @Success 200 {object} models.User
// @Failure 403 id is empty
// @router /getDev/ [get]
//func (c *UserController) GetDev() {
//
//	us, err := models.GetDev(c.Uid())
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(us)
//}

// @Title 按部门获取用户
// @router /getbydep/ [get]
func (c *UserController) GetByDep() {
	users := models.GetAllUserByDp()
	c.RespJSONData(users)
}

// @Title 获取商务负责人，以及每个负责人所负责的游戏
// @router /getBusinessPeople/ [get]
func (c *UserController) GetBusinessPeople() {
	us, err := models.GetUsersByDevMent(237) //获取商务部所有用户

	models.BusinessAddGameIdAndChannelCodeInfo(&us)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(us)
}

// @Title 发送邮件
// @router /sendemail/ [post]
func (c *UserController) SendEmail() {
	v := models.GetParam{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	to := v.ToUser
	subject := v.Subject
	body := v.Body
	depId := v.DepId
	toGroup := v.ToGroup
	userType := v.UserType
	if userType == "2" {
		toGroups := strings.Split(toGroup, ";")
		tool.SendMessageEmail(subject, body, toGroups)
		return
	}
	if len(to) != 0 {
		tool.SendMessageEmail(subject, body, [] string{to})
		c.RespJSONData("success")
		return
	}
	tool.SendMessageEmailByDp(subject, body, depId)
	c.RespJSONData("success")
	return

}

// @Description 获取所有人  应用于快递管理页面 下拉选择框
// @router /userList/ [get]
func (c *UserController) GetUserList() {
	users, err := models.GetUserList()
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	c.RespJSONData(users)
}

/*// @Description 根据Department获取所有人
// @router /get_users_for_department/ [get]
func (c *UserController) GetUsersForDepartment(){

	department ,err:= c.GetInt("department")
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	users ,err := models.GetUsersForDepartment(department)
	if err!=nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
	}
	c.RespJSONData(users)
}

// @Description 根据userid获取所有department的人
// @router /get_users_for_userid_by_department/ [get]
func (c *UserController) GetUsersForUserIdByDepartment(){
	id ,err:=c.GetInt("user_id")
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	users , err := models.GetUsersForUserIdByDepartment(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(users)

}*/
