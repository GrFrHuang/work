package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 敏感操作
type UserFaceVerifyController struct {
	BaseController
}

// Post ...
// @Title Post
// @Description create UserFaceVerify
// @Param	body		body 	models.UserFaceVerify	true		"body for UserFaceVerify content"
// @Success 201 {int} models.UserFaceVerify
// @Failure 403 body is empty
// @router / [post]
func (c *UserFaceVerifyController) Post() {
	var v models.UserFaceVerify
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if len(v.Image) == 0 || v.Uid == 0 || v.Tag == 0 || len(v.Code) == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	if err := models.AddUserFaceVerify(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("ok")
	return
}

// Post ...
// @Title Post
// @Description create UserFaceVerify
// @Param	body		body 	models.UserFaceVerify	true		"body for UserFaceVerify content"
// @Success 201 {int} models.UserFaceVerify
// @Failure 403 body is empty
// @router /qrcode [post]
func (c *UserFaceVerifyController) QrCode() {
	var v models.UserFaceVerify
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if v.Tag == 0 || len(v.Code) == 0 || len(v.Token) == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	if err := models.AddUserCodeLogin(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("ok")
	return
}
