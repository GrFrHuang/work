package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strings"
	"fmt"
)

// WarningLogController operations for WarningLog
type WarningLogController struct {
	BaseController
}

// URLMapping ...
func (c *WarningLogController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create WarningLog
// @Param	body		body 	models.WarningLog	true		"body for WarningLog content"
// @Success 201 {int} models.WarningLog
// @Failure 403 body is empty
// @router / [post]
func (c *WarningLogController) Post() {
	var v models.WarningLog
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddWarningLog(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get WarningLog by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.WarningLog
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WarningLogController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetWarningLogById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get WarningLog
// @Param	start_time	query	string	false	"e.g. 2006-01-02"
// @Param	end_time	query	string	false	"e.g. 2006-01-02"
// @Param	channels	query	string	false	"e.g. code1, code2, ..."
// @Param	games    	query	string	false	"e.g. 1, 2, ..."
// @Param	grade    	query	string	false	"e.g. 1"
// @Param	limit	    query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	    query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WarningLog
// @Failure 403
// @router / [get]
func (c *WarningLogController) GetAll() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_WARNING_LOG, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	startTime, _ := c.GetInt("start_time")
	endTime, _ := c.GetInt("end_time")
	channels := c.GetString("channels")
	games := c.GetString("games")
	grade, _ := c.GetInt("grade")
	fmt.Printf("games:%s", games)

	channel_list := []interface{}{}
	for _, v := range strings.Split(channels, ",") {
		if v != "" {
			channel_list = append(channel_list, v)
		}
	}

	game_list := []interface{}{}
	for _, v := range strings.Split(games, ",") {
		if v != "" {
			game_list = append(game_list, v)
		}
	}

	if startTime != 0 {
		filter.Where["create_time__gte"] = []interface{}{startTime}
	}
	if endTime != 0 {
		filter.Where["create_time__lte"] = []interface{}{endTime}
	}
	if grade != 0 {
		filter.Where["grade"] = []interface{}{grade}
	}

	if len(game_list) != 0 {
		filter.Where["game_id__in"] = game_list
	}

	if len(channel_list) != 0 {
		filter.Where["channel_code__in"] = channel_list
	}

	tool.InjectPermissionWhere(where, &filter.Where)

	filter.Sortby = []string{"CreateTime"}
	filter.Order = []string{"desc"}
	ss := []models.WarningLog{}
	total, err := tool.GetAllByFilterWithTotal(new(models.WarningLog), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONDataWithTotal(ss, total)
}

// Put ...
// @Title Put
// @Description update the WarningLog
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.WarningLog	true		"body for WarningLog content"
// @Success 200 {object} models.WarningLog
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WarningLogController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.WarningLog{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWarningLogById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the WarningLog
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WarningLogController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteWarningLog(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
