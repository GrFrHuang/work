package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"time"
	"fmt"
	"kuaifa.com/kuaifa/work-together/cmd/warning_system/checker"
)

// 游戏接入
type GameController struct {
	BaseController
}

// URLMapping ...
func (c *GameController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description 创建一条游戏信息，游戏提测时调用
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 201 {int} models.Game
// @Failure 403 body is empty
// @router / [post]
func (c *GameController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_GAME_PLAN_PUB, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var game models.Game
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &game); err == nil {
		//game.GameId, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		//添加 增加游戏提测人的信息

		//参数校验
		if game.GameName == "" || game.Issue == 0 || game.GameType == 0 || game.Package == "" || game.Budget == "" ||
			game.Star == "" || game.Source == 0 {
			c.RespJSON(bean.CODE_Forbidden, "参数缺少")
			return
		}

		//检验 该游戏是否已存在
		check := models.CheckGameIsExistByGameName(game.GameName, game.Source)
		if check {
			c.RespJSON(bean.CODE_Forbidden, "该游戏已存在!")
			return
		}

		game.UpdateRefUserID = c.Uid()
		game.UpdateRefTime = time.Now().Unix()
		game.SubmitPerson = c.Uid()

		if _, err := models.AddGame(&game); err == nil {
			c.RespJSONData("OK")
		} else {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}
}

// GetOne ...
// @Title Get One
// @Description get Game by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Game
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GameController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetGameById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(v)
	}
}

// @Description 游戏评测和游戏接入页的下载
// @router /download [get]
func (c *GameController) DownLoad() {
	c.Ctx.Input.SetParam("limit", "0")
	c.Ctx.Input.SetParam("offset", "0")

	_, ss, errCode, err := c.getAll(bean.PMSM_GAME)
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	flag := c.GetString("flag")
	if flag == "pc" { //游戏评测页
		c.pcDownload(ss)
	} else if flag == "jr" { //游戏接入页
		c.jrDownload(ss)
	}
}

// @router /gameAll [get]
func (c *GameController) GameAll() {
	gameName := c.GetString("gameName", "")
	if gameName == "" {
		c.RespJSON(bean.CODE_Forbidden, "参数错误")
		return
	}

	gameAll := models.GetGameAll(gameName)
	c.RespJSONData(gameAll)
}

//游戏测评下载
func (c *GameController) pcDownload(ss []models.Game) {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{}

		//if v.Game != nil {
		//	r["game"] = v.Game.GameName
		//}
		r["game"] = v.GameName
		evaltime := time.Unix(v.UpdateEvalTime, 0).Format("2006-01-02")
		time := time.Unix(v.ResultTime, 0).Format("2006-01-02")
		r["time"] = time
		r["evaltime"] = evaltime
		if v.User != nil {
			r["user"] = v.User.Nickname
		}
		if len(*v.Cepingpeoples) > 0 {
			var peoples []string
			for _, people := range *v.Cepingpeoples {
				peoples = append(peoples, people.Nickname)
			}
			r["user"] = strings.Join(peoples, ",")
		}
		if v.Ceping != nil {
			r["result"] = v.Ceping.Name
		}
		//获取更新人nickname
		if v.UpdateEvalUserID > 0 {
			user, _ := models.GetUserById(v.UpdateEvalUserID)
			r["update_evaluser"] = user.Nickname
		}
		rs[i] = r
	}

	cols := []string{"game", "time", "user", "result", "update_evaluser", "evaltime"}
	maps := map[string]string{
		cols[0]: "游戏名称",
		cols[1]: "评测时间",
		cols[2]: "评测人",
		cols[3]: "评测结果",
		cols[4]: "更新人",
		cols[5]: "更新时间",
	}
	tmpFileName := fmt.Sprintf("游戏评测-%s", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

//游戏接入下载
func (c *GameController) jrDownload(ss []models.Game) {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{
			"game_name":       v.GameName,
			"quantity_policy": v.QuantityPolicy,
			"remarks":         v.Remarks,
		}
		var import_time = ""
		var publish_time = ""
		var jrtime = ""
		if v.ImportTime != 0 {
			import_time = time.Unix(v.ImportTime, 0).Format("2006-01-02")
		}
		if v.PublishTime != 0 {
			publish_time = time.Unix(v.PublishTime, 0).Format("2006-01-02")
		}
		if v.UpdateJrTime != 0 {
			jrtime = time.Unix(v.UpdateJrTime, 0).Format("2006-01-02")
		}
		r["import_time"] = import_time
		r["publish_time"] = publish_time
		r["jrtime"] = jrtime
		if v.Leixing != nil {
			r["type"] = v.Leixing.Name
		}
		if v.Hezuo != nil {
			r["cooperation"] = v.Hezuo.Name
		}
		if v.Yanfa != nil {
			r["deve"] = v.Yanfa.Name
		}
		if v.Faxing != nil {
			r["issue"] = v.Faxing.Name
		}
		if v.Gonghui != nil {
			r["sociaty"] = v.Gonghui.Name
		}
		if v.User != nil {
			r["import_people"] = v.User.Nickname
		}
		//获取更新人nickname
		user, _ := models.GetUserById(v.UpdateJrUserID)
		if user != nil {
			r["update_jruser"] = user.Nickname
		}
		rs[i] = r
		if v.GameId > 0 {
			r["state"] = "已接入"
		} else if v.GameId == -1 {
			r["state"] = "不接入"
		} else {
			r["state"] = "未接入"
		}
	}

	cols := []string{"game_name", "import_time", "publish_time", "type", "cooperation", "deve",
		"issue", "quantity_policy", "sociaty", "remarks", "import_people", "update_jruser", "jrtime", "state"}
	maps := map[string]string{
		cols[0]:  "游戏名称",
		cols[1]:  "接入时间",
		cols[2]:  "首发时间",
		cols[3]:  "类型",
		cols[4]:  "合作方式",
		cols[5]:  "研发商",
		cols[6]:  "发行商",
		cols[7]:  "保量政策",
		cols[8]:  "公会政策",
		cols[9]:  "备注",
		cols[10]: "接入人",
		cols[11]: "更新人",
		cols[12]: "更新时间",
		cols[13]: "状态",
	}
	tmpFileName := fmt.Sprintf("游戏接入-%s", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

// GetAll ...
// @Title Get All
// @Description get Game
// @Param	pmsm		query 	string	false		"游戏鉴权类型"
// @Param	gameids		query 	string	false		"游戏id"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Game
// @Failure 403
// @router / [get]
func (c *GameController) GetAll() {
	pmsm_game := c.GetString("pmsm", bean.PMSM_GAME)
	if pmsm_game != bean.PMSM_GAME && pmsm_game != bean.PMSM_CONTRACT_CP && pmsm_game != bean.PMSM_CONTRACT_CHANNEL {
		pmsm_game = bean.PMSM_GAME
	}
	total, ss, errCode, err := c.getAll(pmsm_game)
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

// GetAll ...
// @Title Get Select
// @Description 获取各页面游戏名称下拉列表游戏名及id
// @Param	flag	query	string	false	"标识游戏评测页列表还是游戏提测页列表"
// @Success 200 {object} models.Game
// @Failure 403
// @router /select [get]
func (c *GameController) GetSelect() {

	flag := c.GetString("flag")
	var where map[string][]interface{}
	var err error

	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if flag == "cpht" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CP, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "qdht" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CHANNEL, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "tc" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_PUB, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "pc" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_RESULT, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "jr" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
		filter.Where["result__gt"] = []interface{}{0}
	} else if flag == "sx" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "sxyy" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_OPERATOR, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "sxkf" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_CUSTOMER, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "yxgx" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_UPDATE, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == "qdjr" {
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_CHANNEL_ACCESS, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	}

	//where, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME, nil)
	//if err != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err.Error())
	//	return
	//}

	filter.Fields = append(filter.Fields, "id", "game_id", "game_name")
	filter.Order = append(filter.Order, "desc")
	filter.Sortby = append(filter.Sortby, "create_time")

	if flag == "sx" || flag == "sxyy" || flag == "sxkf" || flag == "cpht" || flag == "qdht" || flag == "yxgx" ||
		flag == "qdjr" {
		filter.Where["game_id__gt"] = []interface{}{0}
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	var ss []models.Game

	err = tool.GetAllByFilter(new(models.Game), &ss, filter)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(ss)
}

// GetAll ...
// @Title Get Select
// @Description 获取添加游戏更新页游戏列表，排除已有游戏更新数据
// @Success 200 {object} models.Game
// @Failure 403
// @router /select/gameUpdate/ [get]
func (c *GameController) GetGameUpdateSelect() {

	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_UPDATE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	games, err := models.GetGameUpdateSelect()

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(games)
}

// GetAll ...
// @Title Get All
// @Description 获取游戏评测列表
// @Param	ids		query 	string	false		"游戏记录的id"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	flag	query	string	false	"标识游戏评测页列表还是游戏提测页列表"
// @Success 200 {object} models.Game
// @Failure 403
// @router /result [get]
func (c *GameController) GetResult() {
	total, ss, errCode, err := c.getAll(bean.PMSM_GAME)
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

// @Description 获取最近新增的游戏
// @router /latest [get]
func (c *GameController) Getlatest() {
	games, err := models.GetLatestGame()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(games)
	}
}

func (c *GameController) getAll(pmsm_game string) (total int64, ss []models.Game, errCode int, err error) {
	flag := c.GetString("flag")
	if flag == "tjqd" {
		pmsm_game = bean.PMSM_GAME_CHANNEL_ACCESS
	}

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, pmsm_game, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	//filter.Order = append(filter.Order, "desc")
	//filter.Sortby = append(filter.Sortby, "publish_time")
	filter.Fields = append(filter.Fields, "id", "game_id", "game_name")

	filter.Where = map[string][]interface{}{}
	if flag == "pc" { //游戏评测页
		//filter.Where["game_id__exact"] = []interface{}{0}
		filter.Fields = append(filter.Fields, "result_time", "result_person", "result", "result_report_word",
			"result_report_excel", "package", "advise", "update_evaluserid", "update_evaltime",
			"submit_person")
	} else if flag == "jr" { //游戏接入页
		//filter.Where["game_id__exact"] = []interface{}{0}
		filter.Where["result__gt"] = []interface{}{0}
		filter.Fields = append(filter.Fields, "import_time", "publish_time", "game_type", "cooperation",
			"development", "issue", "quantity_policy", "sociaty_policy", "remarks", "access_person", "ladders",
			"update_jruserid", "update_jrtime", "submit_person", "body_my", "source")
	} else if flag == "tc" { //游戏提测页
		//filter.Where["game_id__exact"] = []interface{}{0}
		filter.Fields = append(filter.Fields, "publish_time", "issue", "game_type", "ip", "star", "budget",
			"picture", "package", "remarks", "number", "update_refuserid", "update_reftime", "create_time",
			"submit_person", "body_my", "source")
	} else if flag == "tjqd" { //添加渠道
		filter.Fields = append(filter.Fields, "publish_time")
	}

	idStr := c.GetString("ids")
	var ids []interface{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}
	if len(ids) != 0 {
		filter.Where["id__in"] = ids
	}

	peopleStr := c.GetString("peoples")
	if peopleStr != "" {
		var peoples []interface{}
		peoples = append(peoples, peopleStr)
		if len(peoples) != 0 {
			filter.Where["result_person__contains"] = peoples
		}
	}

	resultStr := c.GetString("result")
	if resultStr != "" {
		var result []interface{}
		result = append(result, resultStr)
		if len(result) != 0 {
			filter.Where["result__exact"] = result
		}
	}

	timeStr := c.GetString("time")
	if timeStr != "" {
		var resultTime []interface{}
		for _, v := range strings.Split(timeStr, ",") {
			if v != "" {
				//tm, _ := time.Parse("2006-01-02", v)
				//resultTime = append(resultTime, tm.Unix())
				resultTime = append(resultTime, v)
			}
		}
		if flag == "jr" || flag == "tc" { //提测 接入 首发时间 条件排序
			if len(resultTime) != 0 {
				filter.Where["publish_time__gte"] = []interface{}{resultTime[0]}
				filter.Where["publish_time__lte"] = []interface{}{resultTime[1]}
			}
		} else {
			if len(resultTime) != 0 {
				filter.Where["result_time__gte"] = []interface{}{resultTime[0]}
				filter.Where["result_time__lte"] = []interface{}{resultTime[1]}
			}
		}
	}

	// 提测人筛选
	submit := c.GetString("submit")
	if submit != "" && flag == "tc" {
		var ids []interface{}
		for _, v := range strings.Split(submit, ",") {
			ids = append(ids, v)
		}
		if len(ids) != 0 {
			filter.Where["submit_person__in"] = ids
		}
	}
	// 提测时间筛选
	subTime := c.GetString("subTime")
	if subTime != "" && flag == "tc" {
		var resultTime []interface{}
		for _, v := range strings.Split(subTime, ",") {
			if v != "" {
				resultTime = append(resultTime, v)
			}
		}
		filter.Where["create_time__gte"] = []interface{}{resultTime[0]}
		filter.Where["create_time__lte"] = []interface{}{resultTime[1]}
	}

	// 接入人筛选
	access := c.GetString("access")
	if access != "" && flag == "jr" {
		var ids []interface{}
		for _, v := range strings.Split(access, ",") {
			ids = append(ids, v)
		}
		if len(ids) != 0 {
			filter.Where["access_person__in"] = ids
		}
	}
	// 接入时间筛选
	accTime := c.GetString("accTime")
	if accTime != "" && flag == "jr" {
		var resultTime []interface{}
		for _, v := range strings.Split(accTime, ",") {
			if v != "" {
				resultTime = append(resultTime, v)
			}
		}
		filter.Where["import_time__gte"] = []interface{}{resultTime[0]}
		filter.Where["import_time__lte"] = []interface{}{resultTime[1]}
	}
	// 接入发行商筛选
	issue := c.GetString("issue", "")
	if issue != "" && flag == "jr" {
		var ids []interface{}
		for _, v := range strings.Split(issue, ",") {
			ids = append(ids, v)
		}
		filter.Where["issue__in"] = ids
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.Game{}
	total, err = tool.GetAllByFilterWithTotal(new(models.Game), &ss, filter)

	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	if flag == "pc" { //游戏评测页
		filter.Where["game_id__exact"] = []interface{}{0}
		filter.Fields = append(filter.Fields, "result_time", "result_person", "result", "result_report_word",
			"result_report_excel")
		if len(peopleStr) != 0 {
			var dd []models.Game
			for _, s := range ss {
				var id []int
				json.Unmarshal([]byte(s.ResultPerson), &id)
				f := false
				for _, i := range id {
					k, _ := strconv.Atoi(peopleStr)
					f = k == i
					if f {
						dd = append(dd, s)
						break
					}
				}
			}
			ss = dd
		}
		models.AddResultCepingInfo(&ss)
		models.AddResultCepingPeopleInfo(&ss)
		models.AddAdviseInfo(&ss)
		models.AddUpdateInfo(&ss, "pc")
		models.AddSubmitPeopleInfo(&ss)
	} else if flag == "jr" { //游戏接入页
		models.AddYanfaInfo(&ss)
		models.AddTypeInfo(&ss)
		models.AddHezuoInfo(&ss)
		models.AddFaxingInfo(&ss)
		models.AddGonghuiInfo(&ss)
		models.AddUserInfo(&ss)
		models.GameLadder2Json(&ss)
		models.AddUpdateInfo(&ss, "jr")
		models.AddAccessStateInfo(&ss)
		models.AddSourceInfo(&ss)
	} else if flag == "tc" {
		models.AddSubmitPeopleInfo(&ss)
		models.AddFaxingInfo(&ss)
		models.AddTypeInfo(&ss)

		models.AddSourceInfo(&ss)
		//models.AddPictureInfo(&ss)
		models.AddUpdateInfo(&ss, "ref")
	} else {
		//models.AddYanfaInfo(&ss)
		//models.AddUpdateInfo(&ss,"ref")
	}

	return
}

// GetTCStatistics ...
// @Title Get All
// @Description 提测统计
// @Success 200 {object} models.Statistics
// @Failure 403
// @router /TCStatistics [get]
func (c *GameController) GetTCStatistics() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_PUB, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	ss := models.GetTCStatistics()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(ss)
	}
}

// GetPCStatistics ...
// @Title Get All
// @Description 评测统计
// @Success 200 {object} models.Statistics
// @Failure 403
// @router /PCStatistics [get]
func (c *GameController) GetPCStatistics() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_PUB, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	ss := models.GetPCStatistics()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(ss)
	}
}

// GetJRStatistics ...
// @Title Get All
// @Description 接入统计
// @Success 200 {object} models.Statistics
// @Failure 403
// @router /JRStatistics [get]
func (c *GameController) GetJRStatistics() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN_PUB, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	ss := models.GetJRStatistics()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(ss)
	}
}

// GetAll ...
// @Title Get All
// @Description 获取要添加的接入游戏列表(在表game_all中，但不在表game中)
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Success 200 {object} models.Game
// @Failure 403
// @router /getAddGame [get]
func (c *GameController) GetAddGame() {
	//l, err := models.GetAllGame(query, fields, sortby, order, offset, limit)
	//where, err2 := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_ALL, []string{"game_id"})
	//if err2 != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err2.Error())
	//	return
	//}
	maps, err := models.GetAddGame()

	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(maps)
	}
}

// Put ...
// @Title Put
// @Description 游戏接入
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 200 {object} models.Game
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GameController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Game{Id: id}
	//flag := c.GetString("flag")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.ImportTime == 0 || v.PublishTime == 0 || v.GameType == 0 || v.Cooperation == 0 || v.Issue == 0 ||
		v.AccessPerson == 0 || v.Ladders == "" || v.BodyMy == 0 || v.Source == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数缺少")
		return
	}

	//游戏接入 更新人添加
	v.UpdateJrUserID = c.Uid()
	v.UpdateJrTime = time.Now().Unix()

	// 2018-01-09 修改,取消自动绑定，手动选择游戏接入并绑定game_id
	//通过gameName获取游戏gameId 实现自动绑定功能
	//v.GameId = models.GetGameIdByGameName(v.GameName)
	//if v.GameId <= 0 {
	//	c.RespJSON(bean.CODE_Forbidden, "没有该游戏记录")
	//	return
	//}

	oldGame, _ := models.GetGameById(id)

	if err := models.GameAccessUpdateGameById(&v, where); err == nil {
		var gameplan models.GamePlan
		gameplan.GameId = v.GameId
		if _, err2 := models.AddGamePlan(&gameplan); err2 != nil {
			c.RespJSON(bean.CODE_Forbidden, "绑定失败")
			return
		}

		var con models.Contract
		con.Ladder = v.Ladders
		con.BodyMy = v.BodyMy
		con.IsMain = 1
		//if _, err1 := models.AddContract(&con, v.GameId, 0, v.Issue, "", c.Uid()); err1 != nil {
		if _, err1 := models.AddContract(&con, v.GameId, 0, "", c.Uid(),""); err1 != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
		//游戏接入报警
		if err = checker.SendNewGameAccessWarning(v.GameId); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

		//游戏首发时间变更报警
		//if(oldGame.PublishTime != v.PublishTime){
		if time.Unix(oldGame.PublishTime, 0).Format("2006-01-02") != time.Unix(v.PublishTime, 0).
			Format("2006-01-02") && oldGame.PublishTime != 0 {
			if err = checker.SendGameReleaseUpdateWarning(v.GameId, oldGame.PublishTime, v.PublishTime); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}
		}
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Put ...
// @Title Put
// @Description 游戏接入修改
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 200 {object} models.Game
// @Failure 403 :id is not int
// @router /gameAccess/:id [put]
func (c *GameController) AccessPut() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Game{Id: id}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.ImportTime == 0 || v.PublishTime == 0 || v.GameType == 0 || v.Cooperation == 0 || v.Issue == 0 ||
		v.Ladders == "" || v.BodyMy == 0 || v.GameId == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数缺少")
		return
	}

	//游戏接入 更新人添加
	v.UpdateJrUserID = c.Uid()
	v.UpdateJrTime = time.Now().Unix()

	// 2018-01-09 修改,取消自动绑定，手动选择游戏接入并绑定game_id
	//通过gameName获取游戏gameId 实现自动绑定功能
	//v.GameId = models.GetGameIdByGameName(v.GameName)
	//if v.GameId <= 0 {
	//	c.RespJSON(bean.CODE_Forbidden, "没有该游戏记录")
	//	return
	//}

	oldGame, _ := models.GetGameById(id)
	distributionCompany, _ := models.GetDistributionCompanyByCompanyId(v.Issue)
	if v.BodyMy == 1 {
		v.AccessPerson = distributionCompany.YunduanResponsiblePerson
	} else {
		v.AccessPerson = distributionCompany.YouliangResponsiblePerson
	}
	if err := models.GameAccessUpdateGameById(&v, where); err == nil {
		var gameplan models.GamePlan
		gameplan.GameId = v.GameId
		if _, err2 := models.AddGamePlan(&gameplan); err2 != nil {
			c.RespJSON(bean.CODE_Forbidden, "绑定失败")
			return
		}

		var con models.Contract
		con.Ladder = v.Ladders
		con.BodyMy = v.BodyMy
		con.IsMain = 1
		if _, err1 := models.AddContract(&con, v.GameId, 0, "", c.Uid(),""); err1 != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
		//游戏接入报警
		if err = checker.SendNewGameAccessWarning(v.GameId); err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

		//游戏首发时间变更报警
		if time.Unix(oldGame.PublishTime, 0).Format("2006-01-02") != time.Unix(v.PublishTime, 0).
			Format("2006-01-02") && oldGame.PublishTime != 0 {
			if err = checker.SendGameReleaseUpdateWarning(v.GameId, oldGame.PublishTime, v.PublishTime); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}
		}
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Put ...
// @Title 游戏接入后对游戏进行修改
// @Description update the Game
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 200 {object} models.Game
// @Failure 403 :id is not int
// @router /GameAccessUpdate/:id [put]
func (c *GameController) GameAccessUpdate() {

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Game{Id: id}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.ImportTime == 0 || v.PublishTime == 0 || v.GameType == 0 || v.Cooperation == 0 || v.Issue == 0 ||
		v.Ladders == "" || v.BodyMy == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数缺少")
		return
	}

	//游戏接入 更新人添加
	v.UpdateJrUserID = c.Uid()
	v.UpdateJrTime = time.Now().Unix()

	oldGame, _ := models.GetGameById(id)

	if err := models.GameAccessUpdateGameById(&v, where); err == nil {
		//游戏首发时间变更报警
		if time.Unix(oldGame.PublishTime, 0).Format("2006-01-02") != time.Unix(v.PublishTime, 0).
			Format("2006-01-02") && oldGame.PublishTime != 0 {
			if err = checker.SendGameReleaseUpdateWarning(v.GameId, oldGame.PublishTime, v.PublishTime); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}
		}

		//修改了我方主体，需要同步cp合同中我方主体
		if v.BodyMy != oldGame.BodyMy {
			if err = models.UpdateCpContractBody(v.BodyMy, oldGame.GameId); err != nil {
				c.RespJSON(bean.CODE_Forbidden, err.Error())
				return
			}
		}
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Put ...
// @Title Put
// @Description 游戏提测修改
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 200 {object} models.Game
// @Failure 403 :id is not int
// @router /reference/:id [put]
func (c *GameController) ReferenceUpdate() {

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_PUB, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Game{Id: id}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//检验数据
	if v.GameName == "" || v.Issue == 0 || v.GameType == 0 || v.Ip == 0 || v.Star == "" || v.Budget == "" ||
		v.Package == "" || v.Source == 0 {
		c.RespJSON(bean.CODE_Forbidden, "参数缺少")
		return
	}

	//检验该游戏是否已存在
	check, err := models.GetGameById(id)
	if check.GameName != v.GameName {
		// 若果游戏名进行了修改，则查询新游戏名是否已存在
		newCheck := models.CheckGameIsExistByGameName(v.GameName, v.Source)
		if newCheck {
			c.RespJSON(bean.CODE_Forbidden, "该游戏已存在!")
			return
		}
	}

	//更新人 当前操作人的信息
	v.UpdateRefUserID = c.Uid()
	v.UpdateRefTime = time.Now().Unix()

	if err := models.ReferenceUpdateGameById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Put ...
// @Title Put
// @Description 游戏测评修改
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Game	true		"body for Game content"
// @Success 200 {object} models.Game
// @Failure 403 :id is not int
// @router /evaluation/:id [put]
func (c *GameController) EvaluationUpdate() {

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_RESULT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Game{Id: id}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//测评检验
	if v.ResultTime == 0 || v.Result == 0 || v.Advise == "" || v.Package == "" {
		c.RespJSON(bean.CODE_Forbidden, "参数出问题")
		return
	}
	//更新人 当前操作人信息以及时间
	v.UpdateEvalUserID = c.Uid()
	v.UpdateEvalTime = time.Now().Unix()

	if err := models.GameEvaluationUpdateGameById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData("保存成功")
}

// Delete ...
// @Title Delete
// @Description delete the Game
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GameController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteGame(id); err == nil {
		c.RespJSONData("OK")
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}
}
