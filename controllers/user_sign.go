package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/utils"
	"time"
	"math/rand"
	"strconv"
)

// 签到功能
type UserSignController struct {
	BaseController
}

// Post ...
// @Title Post
// @Description create UserSign
// @Param	body		body 	models.FaceJson	true		"body for User content"
// @Success 201 {string} ok
// @Failure 403 body is empty
// @router / [post]
func (c *UserSignController) Post() {
	var face models.FaceJson
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &face); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if face.Tag == 0 || face.Uid == 0 || face.Img == "" || face.Code == "" {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	if err := models.IsExitUSerById(face.Uid); err != nil {
		c.RespJSON(bean.CODE_Params_Err, "用户不存在")
		return
	}
	status := models.GetExistFaceByUid(face.Uid)
	if status != 1 {
		c.RespJSON(bean.CODE_Bad_Request, "该用户资料正在认证中或没认证或被拒绝")
		return
	}

	if err := models.AddUserSign(face.Uid, face.Tag, face.Img, face.Code); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("ok")
	return
}

// Put ...
// @Title 审核签到
// @Description update the User
// @Param	uid		path 	string	true		"用户uid"
// @Param	status		query 	int	true		"状态 1已经审核2未审核3拒绝通过"
// @Param	remarks		query 	int	false		"原因"
// @Success 200 {object} ok
// @Failure 403 :id is not int
// @router /:uid [put]
func (c *UserSignController) Put() {
	/*_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_USER_SIGN_EXAMINE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}*/
	idStr := c.Ctx.Input.Param(":uid")
	id, _ := strconv.Atoi(idStr)
	status, _ := c.GetInt("status", 0)
	remarks := c.GetString("remarks", "")
	if id == 0 || status == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	if status == 3 && remarks == "" {
		c.RespJSON(bean.CODE_Params_Err, "拒绝必须输入原因")
		return
	}
	if status > 3 || status < 1 {
		c.RespJSON(bean.CODE_Params_Err, "status参数范围（1~3）")
		return
	}
	if err := models.ChangeSignStatusByUid(id, status, remarks); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData("修改成功")
	return
}

// GetAll ...
// @Title Get All
// @Description get UserSign
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.UserSign
// @Failure 403
// @router / [get]
func (c *UserSignController) GetAll() {
	/*_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER_SIGN_EXAMINE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}*/
	var ss []models.UserSign
	total, err := tool.GetAllWithTotal(c.Controller, new(models.UserSign), &ss, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

// Put ...
// @Title 判断二维码过期
// @Description update the User
// @Param	code		query 	int	false		"原因"
// @Success 200 {object} ok
// @Failure 403 :id is not int
// @router /code [get]
func (c *UserSignController) Code() {
	code := c.GetString("code", "")
	if len(code) <= 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	st, err := utils.Redis.GET("morning_sign")
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if len(st) <= 0 {
		c.RespJSON(bean.CODE_Params_Err, "未获取到缓存信息")
		return
	}
	if code != st {
		c.RespJSON(bean.CODE_Params_Err, "二维码不正确")
		return
	}
	c.RespJSONData("ok")
	return
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
