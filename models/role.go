package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bjson"
)

type Role struct {
	Id            int    `json:"id,omitempty" orm:"column(id);auto"`
	Name          string `json:"name,omitempty" orm:"column(name);size(50);null"`
	PermissionIds string `json:"permission_ids,omitempty" orm:"column(permission_ids);null"`
	Readonly      int `json:"readonly,omitempty" orm:"column(readonly);null"`

	Permissions *[]Permission `orm:"-" json:"permissions,omitempty"`
}

func (t *Role) TableName() string {
	return "role"
}

func init() {
	orm.RegisterModel(new(Role))
}

// AddRole insert a new Role into database and returns
// last inserted Id on success.
func AddRole(m *Role) (id int64, err error) {
	m.Id = 0
	m.Readonly = 2
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Id. Returns error if
// Id doesn't exist
func GetRoleById(id int) (v *Role, err error) {
	o := orm.NewOrm()
	v = &Role{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateRole updates Role by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoleById(m *Role) (err error) {
	o := orm.NewOrm()
	v := Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Readonly == 1 {
		err = errors.New("readonly")
		return
	}
	_, err = o.Update(m)
	return
}

// DeleteRole deletes Role by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRole(id int) (err error) {
	o := orm.NewOrm()
	v := Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Readonly == 1 {
		err = errors.New("readonly")
		return
	}
	_, err = o.Delete(&Role{Id: id})
	return
}

func GroupAddPermissionInfo(us []Role) bool {
	if len(us) == 0 {
		return false
	}
	rsMap := GetAllPermissionMap()
	if rsMap == nil || len(rsMap) == 0 {
		return false
	}
	for i := range us {
		AddPermissionInfo(&us[i], rsMap)
	}
	return true
}

func AddPermissionInfo(v *Role, rsMap map[int]Permission) bool {
	if rsMap == nil || len(rsMap) == 0 {
		return false
	}
	bj, _ := bjson.New([]byte(v.PermissionIds))
	if l := bj.Len(); l != 0 {
		rrs := []Permission{}
		for i := 0; i < l; i++ {
			if r, ok := rsMap[bj.Index(i).Int()]; ok {
				rrs = append(rrs, r)
			}
		}
		v.Permissions = &rrs
	}
	return true
}

func GetAllRoleMap() map[int]Role {
	o := orm.NewOrm()
	qs2 := o.QueryTable(new(Role))
	rs := []Role{}
	rsMap := map[int]Role{}
	_, err := qs2.All(&rs)
	if err != nil {
		return nil
	}
	for _, r := range rs {
		rsMap[r.Id] = r
	}

	return rsMap
}
