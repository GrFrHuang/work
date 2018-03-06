package models

import "github.com/astaxie/beego/orm"

func QueryTable(ptrStructOrTableName interface{}, where map[string][]interface{}) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(ptrStructOrTableName)
	if where != nil {
		for k, v := range where {
			qs.Filter(k, v...)
		}
	}
	return qs
}
