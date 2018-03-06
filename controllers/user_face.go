package controllers

import (
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/utils"
)

//人脸功能
type UserFaceControllers struct {
	BaseController
}

// Post
// @Title 绑定人脸
// @Description bind FaceUser
// @Param	body		body 	models.FaceJson	true		"body for User content"
// @Success 200 {object} ok
// @Failure 403 body is empty
// @router / [post]
func (c *UserFaceControllers) Post() {
	var face models.FaceJson
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &face); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if face.Uid == 0 || face.Img == "" {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	if err := models.IsExitUSerById(face.Uid); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	status := models.GetExistFaceByUid(face.Uid)
	if status == 1 {
		c.RespJSON(bean.CODE_Bad_Request, "该用户已经通过认证")
		return
	} else if status == 2 {
		c.RespJSON(bean.CODE_Bad_Request, "该用户正在认证中")
		return
	}
	//添加图片地址
	as, err := models.NewAsset().SaveAssetBase64Img(face.Img, face.Uid)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	var user models.UserFace
	user.Uid = face.Uid
	user.Status = 2
	user.Path = strconv.FormatInt(as.Id, 10)
	if err := models.AddUserFace(&user); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData("ok")
	return

}

// GET ...
// @Title 创建二维码
// @Description create UserSign
// @Param	tag	query	string	false	"创建二维码"
// @Param	condition	query	string	false	"创建条件"
// @Success 201 {string} ok
// @Failure 403 body is empty
// @router /qrcode [get]
func (c *UserFaceControllers) Create() {
	tag := c.GetString("tag", "")
	if tag == "" {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	var redisKey string
	var v = tag + "@" + GetRandomString(35)
	if tag == "1" {
		redisKey = "morning_sign"
		st, err := utils.Redis.GET(redisKey)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		if len(st) > 0 {
			c.RespJSONData(st)
			return
		}
	} else if tag == "2" {
		redisKey = v
	} else if tag == "3" {
		redisKey = "face_check"
		st, err := utils.Redis.GET(redisKey)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		if len(st) > 0 {
			c.RespJSONData(st)
			return
		}
	} else {
		c.RespJSON(bean.CODE_Params_Err, "tag = 1 or 2")
		return
	}
	if err := utils.Redis.SET(redisKey, v, 60*10); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(v)
	return
}

// Put ...
// @Title 获取用户信息
// @Description update the User
// @Param	uid		path 	string	false		"备注"
// @Success 200 {object} ok
// @Failure 403 :id is not int
// @router /:uid [get]
func (c *UserFaceControllers) GetUserById() {
	idStr := c.Ctx.Input.Param(":uid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	err, path := models.GetFacePathByUid(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	v.FacePathId = path
	c.RespJSONData(v)
	return
}

// Put ...
// @Title 更新用户
// @Description update the User
// @Param	uid		path 	string	true		"用户id"
// @Param	status		formData 	int	true		"状态"
// @Param	remarks		formData 	string	false		"备注"
// @Success 200 {object} ok
// @Failure 403 :id is not int
// @router /:uid [put]
func (c *UserFaceControllers) Put() {
	/*_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_USER_FACE_EXAMINE, nil)
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
	if status > 4 || status < 1 {
		c.RespJSON(bean.CODE_Params_Err, "status参数范围（1~4）")
		return
	}
	if status == 3 && remarks == "" {
		c.RespJSON(bean.CODE_Params_Err, "拒绝必须输入原因")
		return
	}

	if err := models.ChangeStatusByUid(id, status, remarks); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData("修改成功")
	return
}

// GetAll ...
// @Title Get All
// @Description get UserFace
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.UserFace
// @Failure 403
// @router / [get]
func (c *UserFaceControllers) GetAll() {
	/*_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER_FACE_EXAMINE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}*/
	var ss []models.UserFace
	total, err := tool.GetAllWithTotal(c.Controller, new(models.UserFace), &ss, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}
