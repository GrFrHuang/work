package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 我方账户信息
type AccountInformationController struct {
	BaseController
}

// URLMapping ...
func (c *AccountInformationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
}

// GetOne ...
// @Title Get One
// @Description get AccountInformation by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.AccountInformation
// @Failure 403 :id is empty
// @router /:id [get]
//func (c *AccountInformationController) GetOne() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v, err := models.GetAccountInformationById(id)
//	if err != nil {
//		c.Data["json"] = err.Error()
//	} else {
//		c.Data["json"] = v
//	}
//	c.ServeJSON()
//}

// GetAccountInfo ...
// @Title GetAccountInfo
// @Description 根据主体，获取账户信息
// @Param	body_my		query 	string	true		"The key for staticblock"
// @Success 200 {object} models.AccountInformation
// @Failure 403 :id is empty
// @router / [get]
func (c *AccountInformationController) GetAccountInfo() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数：我方主体")
		return
	}

	info, err := models.GetAccountInformation(bodyMy)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(info)
}

// Put ...
// @Title Put
// @Description update the AccountInformation
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.AccountInformation	true		"body for AccountInformation content"
// @Success 200 {object} models.AccountInformation
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AccountInformationController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.AccountInformation{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if err := models.UpdateAccountInformationById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}
