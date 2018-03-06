package controllers

import (
	"errors"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"

	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 所有游戏列表
type GameAllController struct {
	BaseController
}

// URLMapping ...
func (c *GameAllController) URLMapping() {
	//c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("Put", c.Put)
	//c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create GameAll
// @Param	body		body 	models.GameAll	true		"body for GameAll content"
// @Success 201 {int} models.GameAll
// @Failure 403 body is empty
// @router / [post]
//func (c *GameAllController) Post() {
//	var v models.GameAll
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if _, err := models.AddGameAll(&v); err == nil {
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
// @Description get GameAll by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.GameAll
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GameAllController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetGameAllById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get GameAll
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.GameAll
// @Failure 403
// @router / [get]
func (c *GameAllController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 0
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.RespJSON(bean.CODE_Bad_Request, errors.New("Error: invalid query key/value pair"))
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_GAME_ALL, []string{"id"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	l, err := models.GetAllGameAll(query, fields, sortby, order, offset, limit, where)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		c.RespJSONData(l)
	}
}

// GetBundGame ...
// @Title Get Bund game
// @Description 获取绑定游戏game_id及game_name,在game_all表中，不在game表中
// @Success 200 {object} models.GameAll
// @Failure 403
// @router /bund [get]
func (c *GameAllController) GetBundGame(){

	maps, err := models.GetAddGame()
	if err != nil{
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	} else {
		//gamedetail := GetGameDetails(l)
		c.RespJSONData(maps)
	}

	//filter, err := tool.BuildFilter(c.Controller, 0)
	//if err != nil {
	//	c.RespJSON(bean.CODE_Params_Err, err.Error())
	//	//errCode = bean.CODE_Params_Err
	//	return
	//}
	//
	//filter.Fields = append(filter.Fields, "game_id", "game_name")
	//
	//tool.InjectPermissionWhere(nil, &filter.Where)
	//
	//ss := []models.GameAll{}
	//
	//err = tool.GetAllByFilter(new(models.GameAll), &ss, filter)
	////total, err := tool.GetAllByFilterWithTotal(new(models.GameAll), &ss, filter)
	////
	////fmt.Printf("total:%v\n", total)
	//
	//if err != nil {
	//	//errCode = bean.CODE_Not_Found
	//	c.RespJSON(bean.CODE_Not_Found, err.Error())
	//	return
	//}

	//c.RespJSONData(ss)



}

// Put ...
// @Title Put
// @Description update the GameAll
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GameAll	true		"body for GameAll content"
// @Success 200 {object} models.GameAll
// @Failure 403 :id is not int
// @router /:id [put]
//func (c *GameAllController) Put() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v := models.GameAll{Id: id}
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if err := models.UpdateGameAllById(&v); err == nil {
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
// @Description delete the GameAll
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
//func (c *GameAllController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.DeleteGameAll(id); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}
