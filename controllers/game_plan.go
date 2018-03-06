package controllers

import (
	"encoding/json"
	"errors"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/utils"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"fmt"
	"time"
	"kuaifa.com/kuaifa/work-together/tool"
	"github.com/astaxie/beego"
)

// 上线前准备
type GamePlanController struct {
	BaseController
}

// URLMapping ...
func (c *GamePlanController) URLMapping() {
	//c.Mapping("Post", c.Post)
	//c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create GamePlan
// @Param	body		body 	models.GamePlan	true		"body for GamePlan content"
// @Success 201 {int} models.GamePlan
// @Failure 403 body is empty
// @router / [post]
//func (c *GamePlanController) Post() {
//	var v models.GamePlan
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if _, err := models.AddGamePlan(&v); err == nil {
//			c.Ctx.Output.SetStatus(201)
//			c.Data["json"] = v
//		} else {
//			c.Data["json"] = err.Error()
//		}
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// @Title Post
// @Description 运营准备
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GamePlan	true		"body for GamePlan content"
// @Success 201 {int} models.GamePlan
// @Failure 403 body is empty
// @router /operatorPlan/:id [put]
func (c *GamePlanController) OperatorPlan() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GamePlan{Id: id}
	var query = make(map[string]string)

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_OPERATOR, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v.OperatorTime = time.Now().Unix()
		v.OperatorUpdate = c.Uid()
		if err := models.OperateUpdateGamePlanById(&v, where); err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}

		gifts := strings.Split(v.Gifts, ",")

		query["game_id"] = strconv.Itoa(v.GameId)
		oldChannels, err := models.GetAllGamePlanChannel(query, nil, nil, nil, 0, 0)
		maps := utils.ObjListToMapList(oldChannels)

		if err == nil {
			for _, bChannel := range maps {
				channelCode := bChannel["ChannelCode"].(string)
				conn, err := models.GetGamePlanChannel(v.GameId, channelCode)
				if err != nil {
					c.RespJSON(bean.CODE_Params_Err, err.Error())
					return
				}

				conn.GiftStatus = 0
				for _, gift := range gifts {
					if channelCode == gift {
						conn.GiftStatus = 1
					}
				}

				conn.Id = models.GetGamePlanChannelId(v.GameId, channelCode)

				if err := models.UpdateGamePlanChannel(&conn); err != nil {
					c.RespJSON(bean.CODE_Forbidden, err.Error())
					return
				}
			}
			c.RespJSONData("ok")
			return
		} else {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
}

//OperatorPlan备份
/*
func (c *GamePlanController) OperatorPlan() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GamePlan{Id: id}
	var query = make(map[string]string)

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_OPERATOR, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v.OperatorTime = time.Now().Unix()
		v.OperatorUpdate = c.Uid()

		//fmt.Printf("v:%v----\n", v)
		//fmt.Printf("v:%v----\n", v.OperatorPerson)
		if err := models.UpdateGamePlanById(&v, where); err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		newChannels := strings.Split(v.Channels, ",")
		gifts := strings.Split(v.Gifts, ",")

		query["game_id"] = strconv.Itoa(v.GameId)
		oldChannels, _ := models.GetAllGamePlanChannel(query, nil, nil, nil,0 ,0)

		if oldChannels == nil{
			for _, nChannel := range newChannels {
				var con models.Contract
				var conn models.GamePlanChannel
				i, _ := strconv.Atoi(nChannel)
				if _, err1 := models.AddContract(&con, v.GameId, 1, i, c.Uid()); err1 != nil {
					c.RespJSON(bean.CODE_Forbidden, err.Error())
				}
				for _, gift := range gifts{
					conn.GiftStatus = 0
					j, _ := strconv.Atoi(gift)
					if i == j{
						conn.GiftStatus = 1
					}
				}
				if _, err1 := models.AddGamePlanChannel(&conn, v.GameId, i); err1 != nil {
					c.RespJSON(bean.CODE_Forbidden, err.Error())
				}
			}

		}else{
			maps := utils.ObjListToMapList(oldChannels)
			var del []string
			var both []string
			var new []string
			for _, oChannel := range maps{
				var state = 0
				for _, nChannel := range newChannels {
					if strconv.Itoa(oChannel["ChannelId"].(int)) == nChannel{
						both = append(both, nChannel)
						state = 1
					}
				}
				if state == 0{
					del = append(del, strconv.Itoa(oChannel["ChannelId"].(int)))
				}
			}
			for _, nchannel := range newChannels{
				var state = 0
				for _, b := range both{
					if nchannel == b{
						state = 1
					}
				}
				if state == 0{
					new = append(new, nchannel)
				}
			}
			if len(del) != 0{
				for _, delChannel := range del{
					var con models.Contract
					con.State = 4
					delChannelId, _ := strconv.Atoi(delChannel)
					con.Id = models.GetContractIdByGameAndCompany(v.GameId, 1, delChannelId)
					if err1 := models.UpdateContractById(&con); err1 != nil {
						c.RespJSON(bean.CODE_Params_Err, err.Error())
						return
					}
					delGameChannelId := models.GetGamePlanChannelId(v.GameId, delChannelId)
					if err1 := models.DeleteGamePlanChannel(delGameChannelId); err1 != nil {
						c.RespJSON(bean.CODE_Params_Err, err.Error())
						return
					}
				}
			}

			if len(both) != 0{
				for _, bChannel := range both{
					//var conn models.GamePlanChannel
					i, _ := strconv.Atoi(bChannel)
					conn, err := models.GetGamePlanChannelById(v.GameId, i)
					if err != nil{
						c.RespJSON(bean.CODE_Params_Err, err.Error())
						return
					}
					conn.GiftStatus = 0

					for _, gift := range gifts{
						j, _ := strconv.Atoi(gift)
						if i == j{
							conn.GiftStatus = 1
						}
					}
					//id := models.GetGamePlanChannelById(v.GameId, i)
					//conn.Id = id
					//fmt.Printf("id:%v\n", id)
					models.UpdateGamePlanChannel(&conn)
				}
			}

			if len(new) != 0{
				for _, nChannel := range new {
					var con models.Contract
					var conn models.GamePlanChannel
					i, _ := strconv.Atoi(nChannel)
					if _, err1 := models.AddContract(&con, v.GameId, 1, i, c.Uid()); err1 != nil {
						c.RespJSON(bean.CODE_Params_Err, err.Error())
						return
					}
					for _, gift := range gifts{
						j, _ := strconv.Atoi(gift)
						if i == j{
							conn.GiftStatus = 1
						}
					}
					if _, err1 := models.AddGamePlanChannel(&conn, v.GameId, i); err1 != nil {
						c.RespJSON(bean.CODE_Params_Err, err.Error())
						return
					}
				}
			}
			c.RespJSONData("ok")
			return
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
}
*/

// Post ...
// @Title Post
// @Description 客服准备
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GamePlan	true		"body for GamePlan content"
// @Success 201 {int} models.GamePlan
// @Failure 403 body is empty
// @router /customerPlan/:id [put]
func (c *GamePlanController) CustomerPlan() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GamePlan{Id: id}
	var query = make(map[string]string)

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_CUSTOMER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v.CustomerPerson = c.Uid()
		v.CustomerTime = time.Now().Unix()
		if err := models.CustumerUpdateGamePlanById(&v, where); err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		materials := strings.Split(v.Materials, ",")
		packages := strings.Split(v.Packages, ",")
		tests := strings.Split(v.Tests, ",")

		query["game_id"] = strconv.Itoa(v.GameId)
		oldChannels, err := models.GetAllGamePlanChannel(query, nil, nil, nil,0 ,0)
		maps := utils.ObjListToMapList(oldChannels)

		if err == nil{
			for _, bChannel := range maps{
				//var conn models.GamePlanChannel
				channelCode := bChannel["ChannelCode"].(string)
				conn, err := models.GetGamePlanChannel(v.GameId, channelCode)
				if err != nil{
					c.RespJSON(bean.CODE_Params_Err, err.Error())
					return
				}
				conn.Material = 0
				for _, material := range materials{
					if channelCode == material{
						conn.Material = 1
					}
				}

				conn.PackageStatus = 0
				for _, packageStatus := range packages{
					if channelCode == packageStatus{
						conn.PackageStatus = 1
					}
				}

				conn.Test = 0
				for _, test := range tests{
					if channelCode == test{
						conn.Test = 1
					}
				}

				conn.Id = models.GetGamePlanChannelId(v.GameId, channelCode)

				if err := models.UpdateGamePlanChannel(&conn); err != nil{
					c.RespJSON(bean.CODE_Forbidden, err.Error())
					return
				}
			}
			c.RespJSONData("ok")
			return
		}else {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
}

// GetAll ...
// @Title Get All
// @Description get GamePlan
// @Param	gameids		query 	string	false		"游戏id"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	flag	query	string	false	"标识查询类型. 0：上线前准备概况；1：上线前准备游戏评测；2：上线前准备运营；3：上线前准备客服"
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router / [get]
func (c *GamePlanController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss,total)
}


func (c *GamePlanController) getAll()(total int64, ss []models.GamePlan, errCode int, err error){
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN, nil)
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
	flag,err := c.GetInt("flag")
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

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

	tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.GamePlan{}
	total,err = tool.GetAllByFilterWithTotal(new(models.GamePlan), &ss, filter)

	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	models.AddGameInfo(&ss)

	if flag == 0{//概况
		models.AddExactGamePlanChannel(&ss)
		models.AddSDKInfo(&ss)
		models.AddYunyingInfo(&ss)
		models.AddYunyingPeopleInfo(&ss)
		models.AddCepingInfo(&ss)
	} else if flag == 1{//测评
		models.AddCepingInfo(&ss)
		models.AddCepingPeopleInfo(&ss)
	} else if flag ==2 {//运营
		models.AddExactGamePlanChannel(&ss)
		models.AddYunyingInfo(&ss)
		models.AddYunyingPeopleInfo(&ss)
		models.AddYunyingUpdateInfo(&ss)
	} else if flag == 3{//客服
		models.AddExactGamePlanChannel(&ss)
		models.AddKefuPeopleInfo(&ss)
		models.AddSDKInfo(&ss)
	}
	return
}

// @router /download [get]
func (c *GamePlanController) DownLoad() {
	c.Ctx.Input.SetParam("limit","0")
	c.Ctx.Input.SetParam("offset","0")

	_, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	flag,err := c.GetInt("flag")
	if flag == 0{//概况
		c.allDownload(ss)
	} else if flag == 1{//测评
		c.resultDownload(ss)
	} else if flag ==2 {//运营
		c.operatorDownload(ss)
	} else if flag == 3{//客服
		c.customerDownload(ss)
	}
}

//概况下载
func (c *GamePlanController) allDownload(ss []models.GamePlan)  {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{
			"group":v.Group,
		}
		if v.Game != nil {
			r["game"] = v.Game.GameName
		}
		if v.SDK != nil {
			r["sdk"] = v.SDK.Name
		}
		if v.Details != nil {
			beego.Debug(v.Details)
			/*maps := utils.ObjToMap(v.Details)*/
			r["gifts"] = v.Details.Gifts
			r["materials"] = v.Details.Materials
			r["channels"] = v.Details.Channels
			r["packages"] = v.Details.Packages
		}
		if v.Yunying != nil {
			r["issue"] = v.Yunying.Name
		}
		if v.Yunyingpeoples != nil {
			peoples := []string{}
			for _, people := range *v.Yunyingpeoples {
				peoples = append(peoples, people.Nickname)
			}
			r["operator"] = strings.Join(peoples, "、")
		}
		if v.Ceping != nil {
			r["result"] = v.Ceping.Result
		}
		rs[i] = r
	}

	cols := []string{"game", "group", "sdk", "gifts", "materials", "channels",
		"packages", "issue", "operator", "result"}
	maps := map[string]string{
		cols[0]:"游戏名称",
		cols[1]:"是否拉组",
		cols[2]:"SDK接入情况",
		cols[3]:"礼包情况",
		cols[4]:"素材情况",
		cols[5]:"总渠道数",
		cols[6]:"发包情况",
		cols[7]:"运营方",
		cols[8]:"运营负责人",
		cols[9]:"评测结果",
	}
	tmpFileName := fmt.Sprintf("上线前准备（概况）-%s.xlsx", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

//游戏测评下载
func (c *GamePlanController) resultDownload(ss []models.GamePlan)  {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{}

		if v.Game != nil {
			r["game"] = v.Game.GameName
		}
		//time := time.Unix(v.ResultTime, 0).Format("2006-01-02")
		time := time.Unix(0, 0).Format("2006-01-02")
		r["time"] = time
		if v.User != nil {
			r["user"] = v.User.Nickname
		}
		if v.Ceping != nil {
			r["result"] = v.Ceping.Result
		}
		rs[i] = r
	}

	cols := []string{"game", "time", "user", "result"}
	maps := map[string]string{
		cols[0]:"游戏名称",
		cols[1]:"评测时间",
		cols[2]:"评测人",
		cols[3]:"评测结果",
	}
	tmpFileName := fmt.Sprintf("上线前准备（游戏评测）-%s.xlsx", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

//运营准备下载
func (c *GamePlanController) operatorDownload(ss []models.GamePlan)  {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{}
		if v.Game != nil {
			r["game"] = v.Game.GameName
		}
		if v.Details != nil {
			beego.Debug(v.Details)
			/*maps := utils.ObjToMap(v.Details)*/
			r["channels"] = v.Details.Channels
			r["gifts"] = v.Details.Gifts
		}
		if v.Yunying != nil {
			r["issue"] = v.Yunying.Name
		}
		if v.Yunyingpeoples != nil {
			peoples := []string{}
			for _, people := range *v.Yunyingpeoples {
				peoples = append(peoples, people.Nickname)
			}
			r["operator"] = strings.Join(peoples, "、")
		}
		if v.User2 != nil{
			r["operator_update"] = v.User2.Nickname
		}
		time := time.Unix(v.OperatorTime, 0).Format("2006-01-02")
		r["time"] = time
		rs[i] = r
	}

	cols := []string{"game", "channels", "gifts", "issue", "operator", "operator_update", "time"}
	maps := map[string]string{
		cols[0]:"游戏名称",
		cols[1]:"总渠道数",
		cols[2]:"礼包情况",
		cols[3]:"运营方",
		cols[4]:"运营负责人",
		cols[5]:"运营更新人",
		cols[6]:"更新时间",
	}
	tmpFileName := fmt.Sprintf("上线前准备（运营准备）-%s.xlsx", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

//客服准备下载
func (c *GamePlanController) customerDownload(ss []models.GamePlan)  {
	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{
			"group":v.Group,
		}
		if v.Game != nil {
			r["game"] = v.Game.GameName
		}
		if v.SDK != nil {
			r["sdk"] = v.SDK.Name
		}
		if v.Details != nil {
			beego.Debug(v.Details)
			r["materials"] = v.Details.Materials
			r["packages"] = v.Details.Packages
		}
		if v.User != nil {
			r["customer"] = v.User.Nickname
		}
		time := time.Unix(v.CustomerTime, 0).Format("2006-01-02")
		r["time"] = time
		rs[i] = r
	}

	cols := []string{"game", "group", "sdk", "materials", "packages", "customer", "time"}
	maps := map[string]string{
		cols[0]:"游戏名称",
		cols[1]:"是否拉组",
		cols[2]:"SDK接入情况",
		cols[3]:"素材情况",
		cols[4]:"发包情况",
		cols[5]:"更新人",
		cols[6]:"更新时间",
	}
	tmpFileName := fmt.Sprintf("上线前准备（客服准备）-%s.xlsx", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

// getAllChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有的下发渠道列表
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getAllChannelsByGameId [get]
func (c *GamePlanController) GetAllChannelsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetAllChannelsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// getDetailsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏上线前准备详情
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getDetailsByGameId [get]
func (c *GamePlanController) GetDetailsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetDetailsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// getAllChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有的下发渠道列表中已发礼包渠道列表
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getGiftChannelsByGameId [get]
func (c *GamePlanController) GetGiftChannelsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetGiftChannelsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// getAllChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有的下发渠道列表中已发包渠道列表
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getPackageChannelsByGameId [get]
func (c *GamePlanController) GetPackageChannelsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetPackageChannelsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// GetTestChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有的下发渠道列表中已测包渠道列表
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getTestChannelsByGameId [get]
func (c *GamePlanController) GetTestChannelsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetTestChannelsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// getAllChannelsByGameId ...
// @Title Get All
// @Description 根据游戏ID获取该游戏所有的下发渠道列表中已发素材渠道列表
// @Param	gameid	query	string	false	"gameid to query. "
// @Success 200 {object} models.GamePlan
// @Failure 403
// @router /getMaterialChannelsByGameId [get]
func (c *GamePlanController) GetMaterialChannelsByGameId() {

	game_id := c.GetString("gameid")
	id,_ := strconv.Atoi(game_id)
	l, err := models.GetMaterialChannelsByGameId(id)

	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		c.RespJSONData(l)
	}
}

// @Title Put
// @Description update the GamePlan
// @Param	id		path 	string	true		"The id you want to update"
// @Param	flag		query 	string	true		"The flag that sign result(0)/operator(1)/customer(2)"
// @Param	body		body 	models.GamePlan	true		"body for GamePlan content"
// @Success 200 {object} models.GamePlan
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GamePlanController) Put() {
	flagid := c.Ctx.Input.Param("flag")
	flag, _ := strconv.Atoi(flagid)

	var where map[string][]interface{}
	var err error
	if flag == 0{//评测更新
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_RESULT, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else if flag == 1{//运营更新
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_OPERATOR, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}

	} else if flag == 2{//客服更新
		where, err = models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_GAME_PLAN_CUSTOMER, nil)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else{
		c.RespJSON(bean.CODE_Not_Found, errors.New("flag wrong!"))
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.GamePlan{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if err := models.UpdateGamePlanById(&v,where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description delete the GamePlan
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GamePlanController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteGamePlan(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
