package workflow

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
)

type WorkflowHistory struct {
	Id         int    `orm:"column(id);auto"`
	TaskId     int    `orm:"column(task_id);null`
	NodeId     int    `orm:"column(node_id);null" description:"结点ID"`
	UserId     int    `orm:"column(user_id);null" description:"用户ID"`
	UserName   string `orm:"column(user_name);size(100);null" description:"姓名"`
	WorkflowId int    `orm:"column(workflow_id);null" description:"任务流ID"`
	Remarks    string `orm:"column(remarks);size(255);null" description:"备注（操作信息）"`
	CreateTime int64  `orm:"column(create_time);null" description:"创建时间"`
	UpdateTime int64  `orm:"column(update_time);null" description:"更新时间"`
}

func (t *WorkflowHistory) TableName() string {
	return "workflow_history"
}

func init() {
	orm.RegisterModel(new(WorkflowHistory))
}

// AddWorkflowHistory insert a new WorkflowHistory into database and returns
// last inserted Id on success.
func AddWorkflowHistory(m *WorkflowHistory) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWorkflowHistoryById retrieves WorkflowHistory by Id. Returns error if
// Id doesn't exist
func GetWorkflowHistoryById(id int) (v *WorkflowHistory, err error) {
	o := orm.NewOrm()
	v = &WorkflowHistory{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWorkflowHistory retrieves all WorkflowHistory matches certain condition. Returns empty list if
// no records exist
func GetAllWorkflowHistory(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WorkflowHistory))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []WorkflowHistory
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateWorkflowHistory updates WorkflowHistory by Id and returns error if
// the record to be updated doesn't exist
func UpdateWorkflowHistoryById(m *WorkflowHistory) (err error) {
	o := orm.NewOrm()
	v := WorkflowHistory{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWorkflowHistory deletes WorkflowHistory by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWorkflowHistory(id int) (err error) {
	o := orm.NewOrm()
	v := WorkflowHistory{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WorkflowHistory{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAll() ([]WorkflowHistory, error) {
	o := orm.NewOrm()
	sql := "select * from workflow_history"
	logs := []WorkflowHistory{}
	_, err := o.Raw(sql).QueryRows(&logs)
	return logs, err
}

func GetLog(where map[string]interface{}) (history []WorkflowHistory, err error) {

	o := orm.NewOrm()
	sql := "select * from workflow_history where 1=1"
	//history := []WorkflowHistory{}
	var valu []interface{}
	for k, v := range where {
		sql = sql + k
		valu = append(valu, v)
	}
	_, err = o.Raw(sql, valu).QueryRows(&history)
	return
}

func AddLog(userId int, task_id int, createTime int64, workflowId int, remarks string, nextName string, flag int) error {


	username, err := models.GetUserById(userId)
	if err != nil {
		return err
	}
	err, nodeid := models.GetDepartmentIdById(userId)
	if err != nil {
		return err
	}
	if flag == 1 {
		remarks = "创建了一个" + remarks
	} else if flag == 2 {
		remarks = "转交任务到" + nextName + "\n 备注：" + remarks
	} else if flag == 3 {
		remarks = "已经完成任务，但还有办理人("+nextName+")未办理\n 备注：" + remarks
	}
	log := WorkflowHistory{
		NodeId:     nodeid,
		TaskId:     task_id,
		UserId:     userId,
		UserName:   username.Nickname,
		CreateTime: createTime,
		WorkflowId: workflowId,
		Remarks:    remarks,
	}
	_, err = AddWorkflowHistory(&log)

	return err

}
