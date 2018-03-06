package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/codes"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bjson"
	"strconv"
)

type CpsServerController struct {
	BaseController
}

// GetOne
// @Title 验证用户接口
// @Description 验证用户接口
// @Success 200 ok
// @Failure 403 :id is empty
// @router /check_user [post]
func (c *CpsServerController) VerifyCheck() {
	role := c.GetString("role", "")
	if role == "" {
		c.RespDataMsg(codes.Params_Err, "参数错误")
		return
	}
	roleId, err := strconv.Atoi(role)
	if err != nil {
		c.RespDataMsg(codes.Params_Err, "参数错误")
		return
	}
	uid := c.Uid()
	user, err := models.GetUserById(uid)
	if err != nil {
		if err != orm.ErrNoRows {
			c.RespDataMsg(codes.Internal_Server_Error, err.Error())
			return
		}
		c.RespDataMsg(codes.Params_Err, "用户不存在")
		return
	}

	flag := false
	//log.Verbose("sb",uid,user)
	bj, err := bjson.New([]byte(user.RoleIds))
	if err != nil {
		c.RespDataMsg(codes.Forbidden, "无权限")
		return
	}
	if l := bj.Len(); l != 0 {
		for i := 0; i < l; i++ {
			if bj.Index(i).Int() == roleId || bj.Index(i).Int() == 1 {
				flag = true
			}
		}
	}
	if !flag {
		c.RespDataMsg(codes.Forbidden, "无权限")
		return
	}

	c.RespDataMsg(codes.Success, "ok")
}

// GetOne
// @Title 获取用户信息
// @Description 获取用户信息
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /get_user [get]
func (c *CpsServerController) GetUserInfo() {
	uid := c.Uid()
	user, err := models.GetUserById(uid)
	if err != nil {
		if err != orm.ErrNoRows {
			c.RespDataMsg(codes.Internal_Server_Error, err.Error())
			return
		}
		c.RespDataMsg(codes.Params_Err, "用户不存在")
		return
	}
	user.RefreshToken = ""
	user.PicFileId = 0
	user.DepartmentId = 0
	user.DepartmentId = 0
	user.RoleIds = ""
	user.Phone = ""
	user.OpenId = ""
	c.RespJSON(codes.Success, user)
	return
}

// @router /get_user_by_department [get]
func (c *CpsServerController) GetUserByDepartment(){

	departmentId := c.GetString("department_id")

	if departmentId == "" {
		c.RespJSON(codes.Params_Err, "参数为空！")
	}
	users , err := models.GetUserByDepartment(departmentId)

	if err != nil {
		c.RespJSON(codes.Params_Err, err.Error())
	}

	c.RespJSONData(users)

}