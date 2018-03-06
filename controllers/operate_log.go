package controllers

import (
	"errors"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
)

// 操作日志
type OperateLogController struct {
	BaseController
}

// URLMapping ...
func (c *OperateLogController) URLMapping() {
}

// GetLogByPage ...
// @Title Get All Log By page
// @Description get Game
// @Param	id	path 	string	true	"The id for contract"
// @Param	page	query 	string	true	"日志页面"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Game
// @Failure 403
// @router /contract/:id [get]
func (c *OperateLogController) GetLogByContractId() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

func (c *OperateLogController) getAll() (total int64, ss []models.OperateLog, errCode int, err error) {

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	page := c.GetString("page")
	if idStr == "" || page == "" {
		errCode = bean.CODE_Bad_Request
		err = errors.New("request parameter missing")
		return
	}
	id, _ := strconv.Atoi(idStr)

	filter.Where["page__exact"] = []interface{}{page}
	filter.Where["data_id__exact"] = []interface{}{strconv.Itoa(id)}
	filter.Order = append(filter.Order, "desc")
	filter.Sortby = append(filter.Sortby, "create_time")

	//tool.InjectPermissionWhere(where, &filter.Where)

	ss = []models.OperateLog{}
	total, err = tool.GetAllByFilterWithTotal(new(models.OperateLog), &ss, filter)

	models.AddOperatePeopleInfo(&ss)

	return
}

// @Title 获取page页面的date_id的资源的action类型日志
// @Description get 获取page页面的data_id的资源的action类型日志(query可传字段:page,data_id,action)
// @Param	query	query 	string	true	"query"
// @Param	limit	query	string	false	"Limit"
// @Param	offset	query	string	false	"Offset"
// @Success 200 {object} models.OperateLog
// @Failure 403
// @router / [get]
func (c *OperateLogController) GetLogs() {
	data := []models.OperateLog{}
	total, err := tool.GetAllWithTotal(c.Controller, new(models.OperateLog), &data, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	models.AddOperatePeopleInfo(&data)
	c.RespJSONDataWithTotal(data, total)
}
