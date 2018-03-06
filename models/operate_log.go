package models

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
	"reflect"
	"strconv"
	"time"
)

type OperateLog struct {
	Id         int    `orm:"column(id);auto" json:"id,omitempty"`
	DataId	   int    `orm:"column(data_id)" json:"data_id,omitempty"`
	UserId     int    `orm:"column(user_id)" json:"user_id,omitempty"`
	Page       int    `orm:"column(page)" json:"page,omitempty"`
	Action     string `orm:"column(action);size(255)" json:"action,omitempty"`
	Content    string `orm:"column(content)" json:"content,omitempty"`
	CreateTime int64  `orm:"column(create_time)" json:"create_time,omitempty"`

	User *User  `orm:"-" json:"user,omitempty"`
}

func (t *OperateLog) TableName() string {
	return "operate_log"
}

func init() {
	orm.RegisterModel(new(OperateLog))
}

func AddOperatePeopleInfo(ss *[]OperateLog) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UserId
	}
	games := []User{}
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	gameMap := map[int]User{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.UserId]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

// AddOperateLog insert a new OperateLog into database and returns
// last inserted Id on success.
func AddOperateLog(m *OperateLog) (id int64, err error) {
	o := orm.NewOrm()

	m.CreateTime = time.Now().Unix()

	id, err = o.Insert(m)
	return
}

// GetOperateLogById retrieves OperateLog by Id. Returns error if
// Id doesn't exist
func GetOperateLogById(id int) (v *[]OperateLog, err error) {
	o := orm.NewOrm()
	v = &[]OperateLog{}
	qs := o.QueryTable(new(OperateLog))
	qs = qs.Filter("content__contains", "\"id\":"+strconv.Itoa(id)).OrderBy("-create_time")
	if _, err = qs.All(v); err == nil {
		return v, nil
	}
	return nil, err
}

//比较新老两条数据，如果数据不一致，即用户有修改，则生成对应的操作日志
func CompareAndAddOperateLog(old interface{}, new interface{}, userId int, page int, id int, action string, needUpdateFields ...string) (err error) {

	oldValue := reflect.Value{}
	newValue := reflect.Value{}
	var oldTypes reflect.Type
	var newTypes reflect.Type
	oldContent := map[string]interface{}{}
	newContent := map[string]interface{}{}

	if(old == nil && new == nil){
		err = errors.New("the two model can not be both nil")
		return
	}else if(old == nil && new != nil){//新增
		newValue = reflect.Indirect(reflect.ValueOf(new))
		newTypes = newValue.Type()
		fields := newValue.NumField()

		for i := 0; i < fields; i ++ {
			newField := newValue.Field(i)
			fieldName := newTypes.Field(i).Name
			tag := newTypes.Field(i).Tag.Get("orm")

			if tag == "-" {
				if utils.ItemInArray(fieldName, needUpdateFields) {
					newContent[fieldName] = newField.Interface()
				}
				continue
			}

			if needUpdateFields != nil && len(needUpdateFields) != 0 {
				if !utils.ItemInArray(fieldName, needUpdateFields) {
					continue
				}
			}

			newContent[fieldName] = newField.Interface()
		}

		if len(newContent) != 0 {
			log := OperateLog{}
			content := map[string]interface{}{}

			log.UserId = userId
			log.Page = page
			log.Action = action
			log.DataId = id
			content["new"] = newContent
			logContent, _ := json.Marshal(content)
			log.Content = string(logContent)
			log.CreateTime = time.Now().Unix()

			_, err = AddOperateLog(&log)
		}
		return
	}else if(old != nil && new == nil){//删除
		oldValue = reflect.Indirect(reflect.ValueOf(old))
		oldTypes = oldValue.Type()
		fields := oldValue.NumField()

		for i := 0; i < fields; i ++ {
			oldField := oldValue.Field(i)
			fieldName := oldTypes.Field(i).Name
			tag := oldTypes.Field(i).Tag.Get("orm")

			if tag == "-" {
				if utils.ItemInArray(fieldName, needUpdateFields) {
					oldContent[fieldName] = oldField.Interface()
				}
				continue
			}

			if needUpdateFields != nil && len(needUpdateFields) != 0 {
				if !utils.ItemInArray(fieldName, needUpdateFields) {
					continue
				}
			}

			oldContent[fieldName] = oldField.Interface()
		}

		if len(oldContent) != 0 {
			log := OperateLog{}
			content := map[string]interface{}{}

			log.UserId = userId
			log.Page = page
			log.Action = action
			log.DataId = id
			content["old"] = oldContent
			logContent, _ := json.Marshal(content)
			log.Content = string(logContent)
			log.CreateTime = time.Now().Unix()

			_, err = AddOperateLog(&log)
		}
		return
	}else{//修改
		newValue = reflect.Indirect(reflect.ValueOf(new))
		oldValue = reflect.Indirect(reflect.ValueOf(old))

		newTypes = newValue.Type()
		oldTypes = oldValue.Type()

		if oldTypes != newTypes {
			err = errors.New("the type of two data is different")
			return
		}
		fields := newValue.NumField()

		for i := 0; i < fields; i ++ {
			oldField := oldValue.Field(i)
			newField := newValue.Field(i)
			fieldName := oldTypes.Field(i).Name
			tag := oldTypes.Field(i).Tag.Get("orm")

			if tag == "-" {
				if utils.ItemInArray(fieldName, needUpdateFields) {
					if oldField.Interface() != newField.Interface() {
						oldContent[fieldName] = oldField.Interface()
						newContent[fieldName] = newField.Interface()
					}
				}
				continue
			}

			if needUpdateFields != nil && len(needUpdateFields) != 0 {
				if !utils.ItemInArray(fieldName, needUpdateFields) {
					continue
				}
			}

			if oldField.Interface() != newField.Interface() {
				oldContent[fieldName] = oldField.Interface()
				newContent[fieldName] = newField.Interface()
			}
		}
		//有内容更新,则生成操作日志
		if len(oldContent) != 0 && len(newContent) != 0 {

			log := OperateLog{}
			content := map[string]interface{}{}

			log.UserId = userId
			log.Page = page
			log.Action = action
			log.DataId = id
			content["old"] = oldContent
			content["new"] = newContent
			logContent, _ := json.Marshal(content)
			log.Content = string(logContent)
			log.CreateTime = time.Now().Unix()

			_, err = AddOperateLog(&log)
		}
	}

	return
}
