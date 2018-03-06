package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"

	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 渠道
type ChannelController struct {
	BaseController
}

// URLMapping ...
func (c *ChannelController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// GetAll ...
// @Title Get All
// @Description get Development
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Development
// @Failure 403
// @router / [get]
func (c *ChannelController) GetAll() {

	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss,total)








	//var fields []string
	//var sortby []string
	//var order []string
	//var query = make(map[string]string)
	//var limit int64 = 0
	//var offset int64
	//
	//// fields: col1,col2,entity.col3
	//if v := c.GetString("fields"); v != "" {
	//	fields = strings.Split(v, ",")
	//}
	//// limit: 10 (default is 10)
	//if v, err := c.GetInt64("limit"); err == nil {
	//	limit = v
	//}
	//// offset: 0 (default is 0)
	//if v, err := c.GetInt64("offset"); err == nil {
	//	offset = v
	//}
	//// sortby: col1,col2
	//if v := c.GetString("sortby"); v != "" {
	//	sortby = strings.Split(v, ",")
	//}
	//// order: desc,asc
	//if v := c.GetString("order"); v != "" {
	//	order = strings.Split(v, ",")
	//}
	//// query: k:v,k:v
	//if v := c.GetString("query"); v != "" {
	//	for _, cond := range strings.Split(v, ",") {
	//		kv := strings.SplitN(cond, ":", 2)
	//		if len(kv) != 2 {
	//			c.RespJSON(bean.CODE_Bad_Request, errors.New("Error: invalid query key/value pair"))
	//			return
	//		}
	//		k, v := kv[0], kv[1]
	//		query[k] = v
	//	}
	//}
	//
	//l, err := models.GetAllChannel(query, fields, sortby, order, offset, limit)
	//if err != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err.Error())
	//} else {
	//	c.RespJSONData(l)
	//}
}

func (c *ChannelController) getAll()(total int64, ss []models.Channel, errCode int, err error){
	//where := map[string][]interface{}{}
	//where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN, nil)
	//if err != nil {
	//	errCode = bean.CODE_Forbidden
	//	return
	//}

	//filter, err := tool.BuildFilter(c.Controller, 0)
	//if err != nil {
	//	errCode = bean.CODE_Params_Err
	//	return
	//}
	//
	//filter.Fields = append(filter.Fields, "channel_id", "name", "cp")
	//
	//tool.InjectPermissionWhere(where, &filter.Where)
	//
	//ss = []models.Channel{}
	//total,err = tool.GetAllByFilterWithTotal(new(models.Channel), &ss, filter)
	//
	//if err != nil {
	//	errCode = bean.CODE_Not_Found
	//	return
	//}

	//
	//return

	ss, total, err = models.GetAllChannels()
	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	return
}

// GetAddChannel ...
// @Title 根据游戏id，获取该游戏可选渠道列表,排除该游戏已发渠道,和终止合作渠道
// @Param	gameId		path 	string	true		"The key for staticblock"
// @router /add/ [get]
func (c *ChannelController) GetAddChannel(){
	//where := map[string][]interface{}{}
	//where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_PLAN, nil)
	//if err != nil {
	//	errCode = bean.CODE_Forbidden
	//	return
	//}
	gameId, err :=c.GetInt("gameId")
	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	channels, err := models.GetAddChannels(gameId)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}else {
		c.RespJSONData(channels)
	}

}

