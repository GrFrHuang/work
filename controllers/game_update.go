package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"

	"kuaifa.com/kuaifa/work-together/models/bean"
	"fmt"
	"time"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/cmd/warning_system/checker"
)

// 游戏更新
type GameUpdateController struct {
	BaseController
}

// URLMapping ...
func (c *GameUpdateController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	//c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create GameUpdate
// @Param	body		body 	models.GameUpdate	true		"body for GameUpdate content"
// @Success 201 {int} models.GameUpdate
// @Failure 403 body is empty
// @router / [post]
func (c *GameUpdateController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_GAME_UPDATE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var v models.GameUpdate
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		//参数校验
		if v.GameId == 0 || v.GameUpdateTime == 0 || v.UpdateType == 0 || v.Version == "" || v.VersionCode == "" ||
			v.PackageName == "" {
			c.RespJSON(bean.CODE_Forbidden, "参数错误")
			return
		}

		v.CreatePerson = c.Uid()
		v.UpdatePerson = c.Uid()
		v.CreateTime = time.Now().Unix()
		v.UpdateTime = time.Now().Unix()
		if _, err := models.AddGameUpdate(&v); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

		//游戏更新预警
		if err = checker.SendGameUpdateWarning(v.GameId); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}

	c.RespJSONData("保存成功")
}

// GetOne ...
// @Title Get One
// @Description get GameUpdate by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.GameUpdate
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GameUpdateController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetGameUpdateById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(v)
	}
}

// getChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有下发渠道
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getChannelsByGameId [get]
func (c *GameUpdateController) GetChannelsByGameId() {

	game_id := c.GetString("gameid")
	fmt.Printf("gameId:%v\n", game_id)
	id, _ := strconv.Atoi(game_id)
	l, err := models.GetChannelsByGameId(id)

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// GetAll ...
// @Title Get All
// @Description get GameUpdate
// @Param	gameids		query 	string	false		"游戏id"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	flag	query	string	false	"标识查询类型. 0：上线前准备概况；1：上线前准备游戏评测；2：上线前准备运营；3：上线前准备客服"
// @Success 200 {object} models.GameUpdate
// @Failure 403
// @router / [get]
func (c *GameUpdateController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

func (c *GameUpdateController) getAll() (total int64, ss []models.GameUpdate, errCode int, err error) {
	//where := map[string][]interface{}{}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_UPDATE, nil)
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
	filter.Sortby = append(filter.Sortby, "create_time")

	idStr := c.GetString("gameids")

	ids := []interface{}{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}

	filter.Where = map[string][]interface{}{}
	if len(ids) != 0 {
		filter.Where["game_id__in"] = ids
	}

	timeStr := c.GetString("time")
	if timeStr != "" {
		resultTime := []interface{}{}
		for _, v := range strings.Split(timeStr, ",") {
			if v != "" {
				//tm, _ := time.Parse("2006-01-02", v)
				resultTime = append(resultTime, v)
			}
		}
		if len(resultTime) != 0 {
			filter.Where["game_update_time__gte"] = []interface{}{resultTime[0]}
			filter.Where["game_update_time__lte"] = []interface{}{resultTime[1]}
		}
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.GameUpdate{}
	total, err = tool.GetAllByFilterWithTotal(new(models.GameUpdate), &ss, filter)

	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	models.UpdateAddGameInfo(&ss)
	models.UpdateAddCreateUserInfo(&ss)
	models.UpdateAddUpdateUserInfo(&ss)
	models.UpdateAddUpdateChannelInfo(&ss)
	models.UpdateAddNotUpdateChannelInfo(&ss)
	models.UpdateAddStopChannelInfo(&ss)

	return
}

// Put ...
// @Title Put
// @Description update the GameUpdate
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GameUpdate	true		"body for GameUpdate content"
// @Success 200 {object} models.GameUpdate
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GameUpdateController) Put() {
	//where := map[string][]interface{}{}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_UPDATE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GameUpdate{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//参数校验
	if v.GameUpdateTime == 0 || v.UpdateType == 0 || v.Version == "" || v.VersionCode == "" || v.PackageName == "" {
		c.RespJSON(bean.CODE_Forbidden, "参数错误")
		return
	}

	v.UpdatePerson = c.Uid()
	if err := models.UpdateGameUpdateById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

/*// Delete ...
// @Title Delete
// @Description delete the GameUpdate
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GameUpdateController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteGameUpdate(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}*/
