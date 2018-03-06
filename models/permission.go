package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bjson"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
	"strings"
	"sync"
	"time"
)

type Permission struct {
	Id        int    `json:"id,omitempty" orm:"column(id);pk"`
	Name      string `json:"name,omitempty" orm:"column(name);size(50);null"`
	Type      int `json:"type,omitempty" orm:"column(type);null"`
	Model     string `json:"model,omitempty" orm:"column(model);null"`
	Methods   string `json:"methods,omitempty" orm:"column(methods);null"`
	Field     string `json:"field,omitempty" orm:"column(field);null"`
	Condition string `json:"condition,omitempty" orm:"column(condition);null"`
	Readonly  int `json:"readonly,omitempty" orm:"column(readonly);null"`
}

func (t *Permission) TableName() string {
	return "permission"
}

func init() {
	orm.RegisterModel(new(Permission))
}

// AddPermission insert a new Permission into database and returns
// last inserted Id on success.
func AddPermission(m *Permission) (id int64, err error) {
	m.Id = 0
	m.Readonly = 2
	o := orm.NewOrm()
	if o.QueryTable(m).Filter("Name", m.Name).Exist() {
		err = fmt.Errorf("permission named %s is exist", m.Name)
		return
	}
	id, err = o.Insert(m)
	return
}

// GetPermissionById retrieves Permission by Id. Returns error if
// Id doesn't exist
func GetPermissionById(id int) (v *Permission, err error) {
	o := orm.NewOrm()
	v = &Permission{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdatePermission updates Permission by Id and returns error if
// the record to be updated doesn't exist
func UpdatePermissionById(m *Permission) (err error) {
	o := orm.NewOrm()
	v := Permission{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Readonly == 1 {
		err = errors.New("this permission is readonly")
		return
	}
	_, err = o.Update(m)
	return
}

// DeletePermission deletes Permission by Id and returns error if
// the record to be deleted doesn't exist
func DeletePermission(id int) (err error) {
	o := orm.NewOrm()
	v := Permission{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Readonly == 1 {
		err = errors.New("this permission is readonly")
		return
	}
	_, err = o.Delete(&Permission{Id: id})
	return
}

func GetAllPermissionMap() map[int]Permission {
	o := orm.NewOrm()
	qs2 := o.QueryTable(new(Permission))
	rs := []Permission{}
	rsMap := map[int]Permission{}
	_, err := qs2.All(&rs)
	if err != nil {
		return nil
	}
	for _, r := range rs {
		rsMap[r.Id] = r
	}

	return rsMap
}

const (
	Type_supper    = 1
	Type_can       = 2
	Type_condition = 3
	Type_notcan    = 4
)

// 指定要生成哪个表的哪个字段的条件, 若没传入field, 则限制全部的field
// action:1234=>增删改查
func CheckPermission(uid int, action string, model string, fields []string) (where map[string][]interface{}, err error) {
	where = map[string][]interface{}{} // condition => args

	ps, err := GetUserPermissions(uid)
	if err != nil {
		return
	}

	hasModelOpPermission := false   // 是否拥有表的操作权限
	whereMap := map[string]string{} // field => whereString

	for _, v := range ps {
		actionsRule := []string{}
		json.Unmarshal([]byte(v.Methods), &actionsRule)
		if v.Type == Type_supper {
			// 超管, 不限制权限
			return
		} else if v.Type == Type_can {
			// 控制是否 能 操作model
			if v.Model == model && utils.ItemInArray(action, actionsRule) {
				hasModelOpPermission = true
				if v.Condition != "" {
					if _, ok := whereMap[v.Field]; ok {
						// 条件不能在同一字段上
						err = fmt.Errorf("%s only condition is existed", v.Field)
						return
					}
					whereMap[v.Field] = v.Condition
				}
			}
		} else if v.Type == Type_condition {
			// 控制 只能操作指定数据 (生成限制查询条件)
			if model == v.Model {
				if utils.ItemInArray(action, actionsRule) {
					if v.Condition == "" {
						err = errors.New("error condition: " + v.Condition)
						return
					}
					if _, ok := whereMap[v.Field]; ok {
						// 只能条件不能在同一字段上
						// 比如不能这样: 只能查看未支付订单+只能查看已支付订单
						err = fmt.Errorf("%s only condition is existed", v.Field)
						return
					}
					whereMap[v.Field] = v.Condition
				}
			}
		} else if v.Type == Type_notcan {
			// 控制是否 不能 操作model
			if v.Model == model && utils.ItemInArray(action, actionsRule) {
				// 拒绝操作
				err = fmt.Errorf("has not %s model's %s permission", model, action)
				return
			}
		}
	}

	if !hasModelOpPermission {
		err = fmt.Errorf("has not %s model's %s permission", model, action)
		return
	}

	if fields != nil && len(fields) != 0 {
		// 筛选要的字段条件
		for k := range whereMap {
			if !utils.ItemInArray(k, fields) {
				delete(whereMap, k)
			}
		}
	}

	// 转化为beego的orm.filter格式
	// condition : in[x,x]|<x|>x|>=x|<=x|=x|x|!=x|
	for f, cond := range whereMap {
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
			}
			//else if strings.Index(x, "!=") == 0 {
			//	where[f] = []interface{}{x[2:]}
			//	//beego 没有不等的操作符, 暂时不实现
			//}
		}
	}

	return
}

type Menu struct {
	Url     string `json:"url"`
	Methods []string `json:"methods"`
}

func GetCanVisitMenu(uid int) (menus []Menu, isAdmin bool) {
	ps, err := GetUserPermissions(uid)
	if err != nil {
		return
	}
	canOperModels := map[string][]string{}

	for _, v := range ps {
		action := []string{}
		json.Unmarshal([]byte(v.Methods), &action)

		if v.Type == Type_supper {
			// 超管, 不限制权限
			isAdmin = true
			return
		} else if v.Type == Type_can {
			// 控制是否 能 操作model
			canOperModels[v.Model] = action
		} else if v.Type == Type_condition {

		} else if v.Type == Type_notcan {
			// 控制是否 不能 操作model
			can := canOperModels[v.Model]
			if can != nil {
				// 删除不能的
				trueCan := []string{}
				for _, c := range can {
					has := false
					for _, not := range action {
						if c == not {
							has = true
							break
						}
					}
					if !has {
						trueCan = append(trueCan, c)
					}
				}
				canOperModels[v.Model] = trueCan
			}
		}
	}

	menus = []Menu{}
	for model, action := range canOperModels {
		if action == nil || len(action) == 0 {
			delete(canOperModels, model)
			continue
		}
		p := bean.MenuMap[model]

		for _, r := range p.Routers {
			menus = append(menus, Menu{Url: r.Url, Methods: action})
		}
	}

	return
}

type CachePermission struct {
	Ps   *[]Permission
	Time int64
}

var cache map[int]CachePermission = map[int]CachePermission{}
var cacheLocker sync.Mutex
// 由于权限多表查询,放入内存缓存,一分钟更新一次
func GetUserPermissions(uid int, ) (result []Permission, err error) {
	cacheLocker.Lock()
	defer cacheLocker.Unlock()

	if c, ok := cache[uid]; ok {
		//60s有效
		if time.Now().Unix()-c.Time <= 60 {
			result = *c.Ps
			return
		}
	}

	user := User{Id: uid}
	o := orm.NewOrm()
	err = o.QueryTable(new(User)).Filter("Id", user.Id).One(&user, "RoleIds")
	if err != nil {
		return
	}
	// 取得角色ids
	rids := []interface{}{}
	//log.Verbose("sb",uid,user)
	bj, err := bjson.New([]byte(user.RoleIds))
	if err != nil {
		return
	}
	if l := bj.Len(); l != 0 {
		for i := 0; i < l; i++ {
			rids = append(rids, bj.Index(i).Int())
		}
	}
	if len(rids) == 0 {
		err = errors.New("user has not roles")
		return
	}

	// 取得角色
	rs := []Role{}
	_, err = o.QueryTable(new(Role)).Filter("Id__in", rids...).All(&rs, "PermissionIds")
	if err != nil {
		return
	}
	if len(rs) == 0 {
		err = errors.New("user has not roles")
		return
	}

	// 取得权限ids
	permissionIds := []interface{}{}
	for _, v := range rs {
		bj, er := bjson.New([]byte(v.PermissionIds))
		if er != nil {
			err = er
			return
		}
		for i, l := 0, bj.Len(); i < l; i++ {
			permissionIds = append(permissionIds, bj.Index(i).Int())
		}
	}
	if len(permissionIds) == 0 {
		err = errors.New("user roles  has not permission")
		return
	}
	ps := []Permission{}
	_, err = o.QueryTable(new(Permission)).Filter("Id__in", permissionIds...).All(&ps, )
	if err != nil {
		return
	}
	if len(ps) == 0 {
		err = errors.New("user roles  has not permission")
		return
	}

	result = ps
	cache[uid] = CachePermission{
		Ps:   &ps,
		Time: time.Now().Unix(),
	}
	return
}
