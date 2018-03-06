package controllers

import (
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// 查看日志
type LogsController struct {
	BaseController
}

// URLMapping ...
func (c *LogsController) URLMapping() {

}

// GetAll ...
// @Title Get All
// @Description 查看日志
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Logs
// @Failure 403
// @router / [get]
func (c *LogsController) GetAll() {
	us := []models.Logs{}
	f, err := tool.BuildFilter(c.Controller, 40)
	if err != nil {
		return
	}
	f.Order=[]string{"desc"}
	f.Sortby=[]string{"id"}
	total,err := tool.GetAllByFilterWithTotal(new(models.Logs),  &us, f)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	// 添加用户信息
	uids := []interface{}{}
	for _, u := range us {
		uids = append(uids, u.UserId)
	}

	u := &models.User{}
	o := orm.NewOrm()
	uss := []models.User{}
	qs := o.QueryTable(u).Filter("Id__in", uids)
	_,err=qs.All(&uss, "Nickname", "Name","Id")
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	for j, l := range us {
		for i, u := range uss {
			if l.UserId == u.Id {
				us[j].User = &uss[i]
			}
		}
	}

	c.RespJSONDataWithTotal(us,total)
}
