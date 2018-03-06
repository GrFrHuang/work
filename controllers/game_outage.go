package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strings"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"time"
	"kuaifa.com/kuaifa/work-together/task/warning"
	"kuaifa.com/kuaifa/work-together/tool"
	"errors"
	"strconv"
)

// 游戏停运
type GameOutageController struct {
	BaseController
}

// URLMapping ...
func (c *GameOutageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
}

// Post ...
// @Title Post
// @Description create GameOutage
// @Param	body		body 	models.GameOutage	true		"body for GameOutage content"
// @Success 201 {int} models.GameOutage
// @Failure 403 body is empty
// @router / [post]
func (c *GameOutageController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_GAME_OUTAGE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var outages []models.GameOutage
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &outages); err == nil {
		for _, v := range outages {
			// 参数校验
			if v.GameId == 0 || v.IncrTime == 0 || v.ServerTime == 0 || v.RechargeTime == 0 {
				c.RespJSON(bean.CODE_Forbidden, "参数错误")
				return
			}

			v.CreateTime = time.Now().Unix()
			v.CreatePerson = c.Uid()

			if _, err := models.AddGameOutage(&v); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}

			// 修改合同状态
			models.UpdateContractState(v.GameId, v.ServerTime)

			// 报警
			if err = warning.GameOutageWarning(v.GameId); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, "请选择正确的数据!")
	}

	c.RespJSONData("保存成功")
}

// GetOne ...
// @Title Get One
// @Description get ChannelAccess by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.GameOutage
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GameOutageController) GetOne() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_OUTAGE, nil)
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
	v, err := models.GetGameOutageById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// Put ...
// @Title 编辑游戏停运
// @Description update the GameOutage
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GameOutage	true		"body for GameOutage content"
// @Success 200 {object} models.GameOutage
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GameOutageController) Put() {

	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_GAME_OUTAGE, nil)
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
	new := models.GameOutage{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &new); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	old, err := models.GetGameOutageById(id)
	old.Desc = new.Desc
	if err = models.UpdateGameOutageById(old); err != nil {
		c.RespJSON(bean.CODE_Not_Found, "修改失败")
		return
	}

	c.RespJSONData("保存成功")
}

// GetGameName ...
// @Title Get GameName
// @Param	type	query	int true	"which type of game name to get: 1:outage game; 2:all game."
// @router /getGameName
func (c *GameOutageController) GetGameName() {
	typ, err := c.GetInt("type", 0)
	if err != nil || (typ != 1 && typ != 2) {
		c.RespJSON(bean.CODE_Params_Err, errors.New("please choose which type of game to get,must in (1,2)").Error())
		return
	}

	games, err := models.GetOutAgeGameName(typ)

	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(games)
}

// GetAll ...
// @Title Get All
// @Description get GameOutage
// @Param	gameids		query 	string	false		"游戏id"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.GameOutage
// @Failure 403
// @router / [get]
func (c *GameOutageController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

func (c *GameOutageController) getAll() (total int64, ss []models.GameOutage, errCode int, err error) {
	//where := map[string][]interface{}{}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_OUTAGE, nil)
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

	var ids []interface{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}

	filter.Where = map[string][]interface{}{}
	if len(ids) != 0 {
		filter.Where["game_id__in"] = ids
	}

	// 服务器状态筛选
	state, _ := c.GetInt("state", 0)
	if state == 1 {
		// 已关服
		now := time.Now().Unix()

		filter.Where["server_time__lte"] = []interface{}{now}
		filter.Order = []string{"desc"}
		filter.Sortby = []string{"server_time"}
	} else if state == 2 {
		// 未关服
		now := time.Now().Unix()

		filter.Where["server_time__gte"] = []interface{}{now}
		filter.Order = []string{"asc"}
		filter.Sortby = []string{"server_time"}
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.GameOutage{}
	total, err = tool.GetAllByFilterWithTotal(new(models.GameOutage), &ss, filter)

	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	models.GameOutAgeAddGameInfo(&ss)
	models.GameOutAgeAddUpdateUserInfo(&ss)

	return
}
