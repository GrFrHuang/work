package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
}

type Baser interface {
	Query() orm.QuerySeter
	getRowsContainer() interface{}
	ReadOrCreate(col1 string, cols ...string) (int64, error)
	Insert() error
	Read(fields ...string) error
	Update(fields ...string) error
	Delete() error
}

type Base struct {
	Outer Baser `orm:"-" json:"-"`
	//
	Orm orm.Ormer `orm:"-" json:"-"`
	//页码
	Page     int `orm:"-" form:"p"`
	PageSize int `orm:"-"`
}
type BaseQuery struct {
}

type DashboardBasicInfo struct {
	Count    float64    `json:"count" orm:"column(count)"`
	DateTime string    `json:"date_time" orm:"column(date_time)"`
}

type DashboradInfo struct {
	Title string        `json:"title"`
	Tag   string        `json:"tag"`
	Show  bool        `json:"show"`
	Href  string        `json:"href"`
	Pmsm  string
	*DashboardBasicInfo
}

type DashboradLinearInfo struct {
	Title string                    `json:"title"`
	Tag   string                    `json:"tag"`
	Info  []DashboardBasicInfo    `json:"info"`
}

func (m *Base) Norm() orm.Ormer {
	if m.Orm == nil {
		m.Orm = orm.NewOrm()
	}
	return m.Orm
}
func (m *Base) Query() orm.QuerySeter {
	return m.Norm().QueryTable(m.Outer)
}
func (m *Base) ReadOrCreate(col1 string, cols ...string) (int64, error) {
	if _, id, err := m.Norm().ReadOrCreate(m.Outer, col1, cols...); err != nil {
		return id, err
	}
	return 0, nil
}
func (m *Base) InsertOrUpdate(cols ...string) (int64, error) {
	//beego.Warning(m.Outer)
	if id, err := m.Norm().InsertOrUpdate(m.Outer, cols...); err != nil {
		return id, err
	}
	return 0, nil
}

func (m *Base) Insert() error {
	if _, err := m.Norm().Insert(m.Outer); err != nil {
		return err
	}
	return nil
}

func (m *Base) Read(fields ...string) error {
	tOut := m.Outer
	if err := m.Norm().Read(m.Outer, fields...); err != nil {
		m.Outer = tOut
		return err
	}
	return nil
}

func (m *Base) Update(fields ...string) error {
	if _, err := m.Norm().Update(m.Outer, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Base) Delete() error {
	if _, err := m.Norm().Delete(m.Outer); err != nil {
		return err
	}
	return nil
}

func (m *Base) getRowsContainer() interface{} {
	return nil
}

func (m *Base) GetList(page, pageSize uint) interface{} {
	q := m.Query()
	var (
		list  interface{}
		start uint
	)
	list = m.Outer.getRowsContainer()
	if page > 1 {
		start = (page - 1) * pageSize
	} else {
		start = 0
	}
	q.Limit(pageSize, start).All(list)
	return list
}
func (m *Base) GetStartPoint(page, pageSize int) (start int) {
	if page > 1 {
		start = (page - 1) * pageSize
	} else {
		start = 0
	}
	return
}

/*
func migirateTable() {
}
func Syncdb() {
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	migirateTable()
}
*/
