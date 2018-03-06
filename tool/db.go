package tool

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

// 封装一层 getAll 获取列表 ,方便使用
func GetAll(c beego.Controller, model interface{}, contains interface{}, defaultLimit int64) (err error) {
	f, err := BuildFilter(c, defaultLimit)
	if err != nil {
		return
	}
	err = GetAllByFilter(model, contains, f)
	return
}
func GetAllWithTotal(c beego.Controller, model interface{}, contains interface{}, defaultLimit int64) (total int64, err error) {
	f, err := BuildFilter(c, defaultLimit)
	if err != nil {
		return
	}
	total, err = GetAllByFilterWithTotal(model, contains, f)
	return
}

type Filter struct {
	Fields []string
	Sortby []string
	Order  []string
	Where  map[string][]interface{}
	Limit  int64
	Offset int64
}

// 根据输入生成过滤条件
func BuildFilter(c beego.Controller, defaultLimit int64) (f *Filter, err error) {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = defaultLimit
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
		//if limit > 50 {
		//	limit = 50
		//}
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
	// query: k:v;k:v
	if v := c.GetString("query"); v != "" {
		sep := ";"
		for _, cond := range strings.Split(v, sep) {
			kv := strings.SplitN(cond, ":", -1)
			if len(kv) != 2 {
				err = errors.New("invalid query key/value pair")
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	where := map[string][]interface{}{} // condition => args
	// condition : inx,x|<x|>x|>=x|<=x|=x|x|!=x|
	for f, cond := range query {
		// 支持一个字段多个条件
		for _, x := range strings.Split(cond, "|") {
			if strings.Index(x, "in") == 0 {
				args := []interface{}{}
				for _, v := range strings.Split(x[2:], ",") {
					args = append(args, v)
				}
				where[f+"__in"] = args
			} else if strings.Index(x, ">=") == 0 {
				where[f+"__gte"] = []interface{}{x[2:]}
			} else if strings.Index(x, ">") == 0 {
				where[f+"__gt"] = []interface{}{x[1:]}
			} else if strings.Index(x, "<=") == 0 {
				where[f+"__lte"] = []interface{}{x[2:]}
			} else if strings.Index(x, "<") == 0 {
				where[f+"__lt"] = []interface{}{x[1:]}
			} else if strings.Index(x, "=") == 0 {
				where[f ] = []interface{}{x[1:]}
			} else if strings.Index(x, "like") == 0 {
				where[f+"__icontains"] = []interface{}{x[4:]}
			} else {
				where[f ] = []interface{}{x}
			}
			//else if strings.Index(x, "!=") == 0 {
			//	where[f] = []interface{}{x[2:]}
			//	//beego 没有不等的操作符, 暂时不实现
			//}
		}
	}

	f = &Filter{
		Fields: fields,
		Limit:  limit,
		Offset: offset,
		Order:  order,
		Where:  where,
		Sortby: sortby,
	}

	return
}

func GetAllByFilter(model interface{}, contains interface{}, f *Filter) (err error) {
	_, err = GetAllByFilterWithTotal(model, contains, f)
	return
}

func GetAllByFilterWithTotal(model interface{}, contains interface{}, f *Filter) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(model)
	// query k=v

	for k, v := range f.Where {
		if k[0] == '!' {
			qs = qs.Exclude(k[1:], v)
		} else {
			qs = qs.Filter(k, v...)
		}
	}
	total, err = qs.Count()
	if err != nil {
		return
	}

	// order by:
	var sortFields []string
	if len(f.Sortby) != 0 {
		if len(f.Sortby) == len(f.Order) {
			// 1) for each sort field, there is an associated order
			for i, v := range f.Sortby {
				orderby := ""
				if f.Order[i] == "desc" {
					orderby = "-" + v
				} else if f.Order[i] == "asc" {
					orderby = v
				} else {
					err = errors.New("Error: Invalid order. Must be either [asc|desc]")
					return
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(f.Sortby) != len(f.Order) && len(f.Order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range f.Sortby {
				orderby := ""
				if f.Order[0] == "desc" {
					orderby = "-" + v
				} else if f.Order[0] == "asc" {
					orderby = v
				} else {
					err = errors.New("Error: Invalid order. Must be either [asc|desc]")
					return
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(f.Sortby) != len(f.Order) && len(f.Order) != 1 {
			err = errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
			return
		}
	} else {
		if len(f.Order) != 0 {
			err = errors.New("Error: unused 'order' fields")
			return
		}
	}

	qs = qs.OrderBy(sortFields...)

	_, err = qs.Limit(f.Limit, f.Offset).All(contains, f.Fields...)

	return
}
