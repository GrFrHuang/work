package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"time"
)

// MainContractController oprations for MainContract
type MainContractController struct {
	BaseController
}

// URLMapping ...
func (c *MainContractController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post
// @Description create MainContract
// @Param	body		body 	models.MainContract	true		"body for MainContract content"
// @Success 201 {int} models.MainContract
// @Failure 403 body is empty
// @router / [post]
func (c *MainContractController) Post() {
	var v models.MainContract
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// 参数校验
		if v.CompanyType == 0 || v.CompanyId == 0 || v.BodyMy == 0 {
			c.RespJSON(bean.CODE_Forbidden, "参数错误")
			return
		}

		v.UpdatePerson = c.Uid()
		v.EffectiveState = 1
		if _, err := models.AddMainContract(&v); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

		if err = models.AddGameList(v.GameIds, v.CompanyId, v.CompanyType); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}

	c.RespJSONData("保存成功")
}

// GetGameName ...
// @Title GetGameName
// @Description 根据company_id获取该发行商下所有游戏名
// @Param	company_id		query 	string	true		"The key for staticblock"
// @Param	type		query 	string	true		"The key for staticblock"
// @Success 200 {object} models.Game
// @Failure 403 channel_code is empty
// @router /gameName [get]
func (c *MainContractController) GetGameName() {
	company_id, err := c.GetInt("company_id", 0)
	if company_id == 0 || err != nil {
		c.RespJSON(bean.CODE_Bad_Request, "参数错误")
		return
	}

	typ, err := c.GetInt("type", 0)
	if typ != 1 && typ != 2 {
		c.RespJSON(bean.CODE_Bad_Request, "参数错误")
		return
	}

	ss, err := models.GetGameName(company_id, typ)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err)
		return
	}

	c.RespJSONData(ss)
}

// GetOne ...
// @Title Get One
// @Description get MainContract by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.MainContract
// @Failure 403 :id is empty
// @router /:id [get]
func (c *MainContractController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetMainContractById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// Put ...
// @Title Put
// @Description update the MainContract
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.MainContract	true		"body for MainContract content"
// @Success 200 {object} models.MainContract
// @Failure 403 :id is not int
// @router /:id [put]
func (c *MainContractController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.MainContract{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//参数判断
	if v.Id == 0 || v.State == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数缺少")
		return
	}

	if v.State != 149 && v.State != 154 && v.State != 155 {
		if v.BeginTime == "" || v.EndTime == "" {
			c.RespJSON(bean.CODE_Params_Err, "参数缺少")
			return
		}
	}

	v.UpdateTime = time.Now().Unix()
	v.UpdatePerson = c.Uid()

	if err := models.UpdateMainContractById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if err := models.AddGameList(v.GameIds, v.CompanyId, v.CompanyType); err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	c.RespJSONData("OK")
}
