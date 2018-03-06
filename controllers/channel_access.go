package controllers

import (
	"encoding/json"
	"errors"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"time"
)

// 渠道接入
type ChannelAccessController struct {
	BaseController
}

// URLMapping ...
func (c *ChannelAccessController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post
// @Description create ChannelAccess
// @Param	body		body 	models.ChannelAccess	true		"body for ChannelAccess content"
// @Success 201 {int} models.ChannelAccess
// @Failure 403 body is empty
// @router / [post]
func (c *ChannelAccessController) Post() {

	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var v models.ChannelAccess

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//参数校验
		if v.GameId == 0 || len(v.Group) == 0 || v.PublishTime == 0 || v.BodyMy == 0 || v.Cooperation == 0 {
			c.RespJSON(bean.CODE_Forbidden, "参数错误")
			return
		}

		var errStr string
		for _, group := range v.Group {
			v.Ladders = group.Ladder

			for _, channel := range group.Channel {
				//渠道添加人 更新人 信息、时间
				v.UpdateChannelUserID = c.Uid()
				v.UpdateChannelTime = time.Now().Unix()

				v.AccessState = 1
				v.AccessUpdateUser = c.Uid()
				v.AccessUpdateTime = time.Now().Unix()

				v.ChannelCode = channel

				channelCompany, _ := models.GetChannelCompanyByCode(channel)
				if v.BodyMy == 1 {
					v.BusinessPerson = channelCompany.YunduanResponsiblePerson
				} else {
					v.BusinessPerson = channelCompany.YouliangResponsiblePerson
				}

				if _, err1 := models.AddChannelAccess(&v); err1 != nil {
					errStr += err1.Error()
					c.RespJSON(bean.CODE_Forbidden, err1.Error())
				}
				var con models.Contract
				var conn models.GamePlanChannel
				con.Ladder = v.Ladders
				con.BodyMy = v.BodyMy
				con.IsMain = 1
				// 如果渠道为“预充值”，则该合同为“预充值渠道，无需合同”
				if v.Cooperation == 166 {
					con.State = 164
				}

				if _, err1 := models.AddContract(&con, v.GameId, 1, v.ChannelCode, c.Uid(),v.Accessory); err1 != nil {
					c.RespJSON(bean.CODE_Forbidden, err1.Error())
				}
				if _, err1 := models.AddGamePlanChannel(&conn, v.GameId, v.ChannelCode); err1 != nil {
					c.RespJSON(bean.CODE_Forbidden, err1.Error())
				}
			}
		}
		err0 := errors.New(errStr)

		if errStr == "" {
			c.RespJSONData("OK")
		} else {
			c.RespJSON(bean.CODE_Forbidden, err0.Error())
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}
}

// GetOne ...
// @Title Get One
// @Description get ChannelAccess by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ChannelAccess
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ChannelAccessController) GetOne() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数错误")
		return
	}
	v, err := models.GetChannelAccessById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get ChannelAccess
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ChannelAccess
// @Failure 403
// @router / [get]
func (c *ChannelAccessController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)

}

func (c *ChannelAccessController) getAll() (total int64, ss []models.ChannelAccess, errCode int, err error) {
	//where := map[string][]interface{}{}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	filter.Order = append(filter.Order, "desc")
	filter.Sortby = append(filter.Sortby, "update_time")

	filter.Where = map[string][]interface{}{}

	idStr := c.GetString("gameids")
	var ids []interface{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}
	if len(ids) != 0 {
		filter.Where["game_id__in"] = ids
	}

	channeltr := c.GetString("channels")
	var channels []interface{}
	for _, v := range strings.Split(channeltr, ",") {
		if v != "" {
			channels = append(channels, v)
		}
	}
	if len(channels) != 0 {
		filter.Where["channel_code__in"] = channels
	}

	businessStr := c.GetString("business")
	if businessStr != "" {
		var business []interface{}
		business = append(business, businessStr)
		if len(business) != 0 {
			filter.Where["business_person__exact"] = business
		}
	}

	// 渠道接入时间筛选
	accTime := c.GetString("accTime")
	if accTime != "" {
		var resultTime []interface{}
		for _, v := range strings.Split(accTime, ",") {
			if v != "" {
				resultTime = append(resultTime, v)
			}
		}
		filter.Where["create_time__gte"] = []interface{}{resultTime[0]}
		filter.Where["create_time__lte"] = []interface{}{resultTime[1]}
	}

	timeStr := c.GetString("time")
	if timeStr != "" {
		var publishTime []interface{}
		for _, v := range strings.Split(timeStr, ",") {
			if v != "" {
				//tm, _ := time.Parse("2006-01-02", v)
				publishTime = append(publishTime, v)
			}
		}
		if len(publishTime) != 0 {
			filter.Where["publish_time__gte"] = []interface{}{publishTime[0]}
			filter.Where["publish_time__lte"] = []interface{}{publishTime[1]}
		}
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.ChannelAccess{}
	total, err = tool.GetAllByFilterWithTotal(new(models.ChannelAccess), &ss, filter)

	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	models.ChannelAccessAddGameInfo(&ss)
	models.ChannelAccessAddChannelInfo(&ss)
	models.ChannelAccessAddHezuoInfo(&ss)
	models.ChannelAccessAddBusinessInfo(&ss)
	//models.ChannelAccessAddCaiwuInfo(&ss)
	models.ChannelAccessParseLadder2Json(&ss)
	models.ChannelAccessAddUpdateUserInfo(&ss)
	models.ChannelAccessAddStateInfo(&ss)

	return
}

//  ...
// @Title Get All
// @Description 渠道接入统计
// @Success 200 {object} models.Statistics
// @Failure 403
// @router /Statistics [get]
func (c *ChannelAccessController) GetStatistics() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	ss := models.GetChannelAccessStatistics()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(ss)
	}
}

// Put ...
// @Title 编辑渠道接入（商务）
// @Description update the ChannelAccess
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ChannelAccess	true		"body for ChannelAccess content"
// @Success 200 {object} models.ChannelAccess
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ChannelAccessController) Put() {

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数错误")
		return
	}
	v := models.ChannelAccess{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//渠道接入 更新人 添加信息
	v.UpdateChannelUserID = c.Uid()
	v.UpdateChannelTime = time.Now().Unix()

	oldChannelAccess, _ := models.GetChannelAccessById(id)

	if oldChannelAccess.AccessState == 2 {
		c.RespJSON(bean.CODE_Params_Err, "已下架的渠道不能再修改")
		return
	}

	if v.AccessState != oldChannelAccess.AccessState {
		v.AccessUpdateTime = time.Now().Unix()
		v.AccessUpdateUser = c.Uid()
	}
	if err := models.UpdateChannelAccessById(&v, where); err == nil {
		//修改了我方主体，需要同步cp合同中我方主体
		con := models.GetContractByChannelCodeAndGameId(oldChannelAccess.GameId, oldChannelAccess.ChannelCode)
		if v.BodyMy != oldChannelAccess.BodyMy {
			con.BodyMy = v.BodyMy
			//if err = models.UpdateChannelContractBody(v.BodyMy, oldChannelAccess.GameId, oldChannelAccess.ChannelCode); err != nil {
			//	c.RespJSON(bean.CODE_Forbidden, err.Error())
			//	return
			//}
		}
		// 修改了渠道合作方式为“预充值”，则该合同为“预充值渠道，无需合同”
		if v.Cooperation == 166 {
			con.State = 164
		}
		if err = models.UpdateContractById(con); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Put ...
// @Title 审核渠道接入（财务）
// @Description update the ChannelAccess
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ChannelAccess	true		"body for ChannelAccess content"
// @Success 200 {object} models.ChannelAccess
// @Failure 403 :id is not int
// @router /audit/:id [put]
func (c *ChannelAccessController) Audit() {

	//where := map[string][]interface{}{}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.ChannelAccess{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//if (v.State == 3){//如果状态已通过，则生成对应的合同，以及上线准备渠道数据
	//	var con models.Contract
	//	var conn models.GamePlanChannel
	//	con.Ladder = v.Ladders
	//	con.BodyMy = v.BodyMy
	//	if _, err1 := models.AddContract(&con, v.GameId, 1, 0, v.ChannelCode, c.Uid()); err1 != nil {
	//		c.RespJSON(bean.CODE_Forbidden, err1.Error())
	//	}
	//	if _, err1 := models.AddGamePlanChannel(&conn, v.GameId, v.ChannelCode); err1 != nil {
	//		c.RespJSON(bean.CODE_Forbidden, err1.Error())
	//	}
	//}

	if err := models.UpdateChannelAccessById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Delete ...
// @Title Delete
// @Description delete the ChannelAccess
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
//func (c *ChannelAccessController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.DeleteChannelAccess(id); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}
