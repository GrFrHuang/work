package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"net/http"
	"fmt"
	"kuaifa.com/kuaifa/work-together/utils"
)

// 会话，用户登录注销
type SessionController struct {
	BaseController
}

// URLMapping ...
func (c *SessionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title 登录
// @Description 用户登录，创建会话
// @Param	body	body 	models.bean.CreateSession	true		"body for CreateSession content"
// @Success 200 {object} models.bean.OutPutSession
// @Failure 403 body is empty
// @router / [post]
func (c *SessionController) Post() {
	var v bean.CreateSession
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}

	loginRet, err := tool.CreateSession(v)
	if err != nil {
		fmt.Println(err.Error())
		c.RespJSON(http.StatusForbidden, err.Error())
		return
	}

	u, err := models.GetUserInfoByName(loginRet.UserName)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}

	if u.Disabled == 2 {
		c.RespJSON(http.StatusForbidden, "此账号被禁止登陆")
		return
	}

	u.AccessToken = loginRet.AccessToken
	u.RefreshToken = loginRet.RefreshToken
	err = models.UpdateUserToken(u)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}
	status := models.GetExistFaceByUid(u.Id)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.RespJSON(http.StatusOK, bean.OutPutSession{Uid: u.Id, Token: u.AccessToken, Status: status})
}

// Post ...
// @Title 二维码登录
// @Description 用户登录，创建会话
// @Param	code	query	string	false	"二维码"
// @Success 200 {object} models.bean.OutPutSession
// @Failure 403 body is empty
// @router /code_login [post]
func (c *SessionController) CodeLogin() {
	qrCode := c.GetString("code", "")
	if len(qrCode) <= 0 {
		c.RespJSON(http.StatusOK, "参数不正确")
		return
	}
	vcode, err := utils.Redis.GET(qrCode)
	if err != nil {
		c.RespJSON(http.StatusOK, err.Error())
		return
	}
	if len(vcode) <= 0 {
		c.RespJSON(http.StatusOK, "二维码已过期")
		return
	}
	if len(vcode) == 37 {
		c.RespJSON(http.StatusOK, "用户未扫码")
		return
	}

	var v bean.CreateSession
	err = json.Unmarshal([]byte(vcode), &v)
	if err != nil {
		c.RespJSON(http.StatusOK, err.Error())
		return
	}
	u, err := models.GetUserInfoByName(v.Email)
	if err != nil {
		c.RespJSON(http.StatusOK, err.Error())
		return
	}
	if u.Disabled == 2 {
		c.RespJSON(http.StatusOK, "此账号被禁止登陆")
		return
	}
	status := models.GetExistFaceByUid(u.Id)
	if err != nil {
		c.RespJSON(http.StatusOK, err.Error())
		return
	}
	go utils.Redis.DEL(qrCode)
	c.RespJSON(http.StatusOK, bean.OutPutSession{Uid: u.Id, Token: v.Password, Status: status})
}

// Delete ...
// @Title 注销
// @Description 用户注销，删除会话
// @Param	x-token		header 	string	true		"x-token"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SessionController) Delete() {
	token := tool.GetRequestToken(c.Ctx)
	u, err := models.GetUserInfoByToken(token)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = tool.OffLine(u.Name)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.RespJSON(http.StatusOK, http.StatusText(http.StatusOK))
	return
}
