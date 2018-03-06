package controllers

import (
	//"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	//"strconv"

	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// 合同寄出归档统计
type ContractStatisticsController struct {
	BaseController
}

// URLMapping ...
func (c *ContractStatisticsController) URLMapping() {
	//c.Mapping("Post", c.Post)
	//c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("Put", c.Put)
	//c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create ContractStatistics
// @Param	body		body 	models.ContractStatistics	true		"body for ContractStatistics content"
// @Success 201 {int} models.ContractStatistics
// @Failure 403 body is empty
// @router / [post]
//func (c *ContractStatisticsController) Post() {
//	var v models.ContractStatistics
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if _, err := models.AddContractStatistics(&v); err == nil {
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

// GetOne ...
// @Title Get One
// @Description get ContractStatistics by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ContractStatistics
// @Failure 403 :id is empty
// @router /:id [get]
//func (c *ContractStatisticsController) GetOne() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v, err := models.GetContractStatisticsById(id)
//	if err != nil {
//		c.Data["json"] = err.Error()
//	} else {
//		c.Data["json"] = v
//	}
//	c.ServeJSON()
//}

// GetAll ...
// @Title Get All
// @Description get ContractStatistics
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	flag	query	string	false	"标识cp或渠道. Must be cp or qd"
// @Success 200 {object} models.ContractStatistics
// @Failure 403
// @router / [get]
func (c *ContractStatisticsController) GetAll() {
	//var where map[string][]interface{}
	//where, err2 := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CP, nil)
	//if err2 != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err2.Error())
	//	return
	//}
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	flag := c.GetString("flag")
	var typ int

	if flag == "cp" {
		typ = 1
	} else if(flag == "qd"){
		typ = 2
	} else{
		c.RespJSON(bean.CODE_Params_Err, "参数错误")
		return
	}

	ss := []models.ContractStatistics{}
	total, err := models.GetAllContractStatisticWithTotal(&ss, typ, filter.Offset, filter.Limit)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONDataWithTotal(ss, total)

}

// GetAll ...
// @Title Get All
// @Description 查询某一天的合同统计详情
// @Param	time	query	string	false	"时间"
// @Param	flag	query	string	false	"标识cp或渠道. Must be cp or qd"
// @Success 200 {object} models.ContractStatistics
// @Failure 403
// @router /getDayDetail/ [get]
func (c *ContractStatisticsController) GetDayDetail() {
	//var where map[string][]interface{}
	//where, err2 := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CP, nil)
	//if err2 != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err2.Error())
	//	return
	//}
	//filter, err := tool.BuildFilter(c.Controller, 20)
	//if err != nil {
	//	c.RespJSON(bean.CODE_Params_Err, err.Error())
	//	return
	//}

	flag := c.GetString("flag")
	var typ int

	if flag == "cp" {
		typ = 1
	} else if(flag == "qd"){
		typ = 2
	} else{
		c.RespJSON(bean.CODE_Params_Err, "参数错误")
		return
	}

	time := c.GetString("time")
	if time == ""{
		c.RespJSON(bean.CODE_Params_Err, "参数错误")
		return
	}

	ss := []models.ContractStatistics{}
	err := models.GetDetails(&ss, typ, time)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	models.AddPeopleInfo(&ss)

	c.RespJSONData(ss)

}

// Put ...
// @Title Put
// @Description update the ContractStatistics
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ContractStatistics	true		"body for ContractStatistics content"
// @Success 200 {object} models.ContractStatistics
// @Failure 403 :id is not int
// @router /:id [put]
//func (c *ContractStatisticsController) Put() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v := models.ContractStatistics{Id: id}
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if err := models.UpdateContractStatisticsById(&v); err == nil {
//			c.Data["json"] = "OK"
//		} else {
//			c.Data["json"] = err.Error()
//		}
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// Delete ...
// @Title Delete
// @Description delete the ContractStatistics
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
//func (c *ContractStatisticsController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.DeleteContractStatistics(id); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}
