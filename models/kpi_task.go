package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type KpiTask struct {
	Id               int             `orm:"column(id);auto" json:"id"`
	TaskCompleteDate int             `orm:"column(task_complete_date);null" description:"任务完成时间" json:"task_complete_date"`
	TaskPublishDate  int             `orm:"column(task_publish_date);null" description:"发布时间" json:"task_publish_date"`
	AssesseserId     int             `orm:"column(assesseser_id);null" description:"考核人" json:"assesseser_id"`
	PublisherId      int             `orm:"column(publisher_id);null" description:"发布人" json:"publisher_id"`
	PublishState     int8            `orm:"column(publish_state);null" description:"发布状态 0.未发布 1.已发布" json:"publish_state"`
	AuditState       int8            `orm:"column(audit_state);null" description:"审核状态 0.未审核 1.已审核" json:"audit_state"`
	CreateTime       int             `orm:"column(create_time);null" description:"创建时间" json:"create_time"`
	UpdateTime       int             `orm:"column(update_time);null" description:"修改时间" json:"update_time"`
	PublishType      int8            `orm:"column(publish_type);null" description:"发布类型 0.立刻 1.定时" json:"publish_type"`
	DepartmentId     int             `orm:"column(department_id);null" description:"部门" json:"department_id"`
	KpiChildTasks    []*KpiChildTask `orm:"-" json:"kpi_child_tasks"`
	AssesseserName   string          `orm:"-" json:"assesseser_name"`
	PublisherName    string          `orm:"-" json:"publisher_name"`
	TotalScore       float64         `orm:"-" json:"total_score"`
}

func (t *KpiTask) TableName() string {
	return "kpi_task"
}

func init() {
	orm.RegisterModel(new(KpiTask))
}

// AddKpiTask insert a new KpiTask into database and returns
// last inserted Id on success.
func AddKpiTask(m *KpiTask) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	m.CreateTime = int(time.Now().Unix())
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

// GetKpiTaskById retrieves KpiTask by Id. Returns error if
// Id doesn't exist
func GetKpiTaskById(id int) (v *KpiTask, err error) {
	o := orm.NewOrm()
	v = &KpiTask{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateKpiTask updates KpiTask by Id and returns error if
// the record to be updated doesn't exist
func UpdateKpiTaskById(m *KpiTask) (err error) {
	o := orm.NewOrm()
	v := KpiTask{Id: m.Id}
	o.Begin()
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
			//update all child tasks
			err = UpdateChildTasks(m)
			if err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		} else {
			beego.Error(err)
			o.Rollback()
			return
		}
	} else {
		beego.Error(err)
		o.Rollback()
		return

	}
	err = o.Commit()
	if err != nil {
		beego.Error(err)
	}
	return
}

// DeleteKpiTask deletes KpiTask by Id and returns error if
// the record to be deleted doesn't exist
func DeleteKpiTask(id int) (err error) {
	o := orm.NewOrm()
	v := KpiTask{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&KpiTask{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Update kpi task fields
func UpdateKpiTaskFields(task *KpiTask, field ... string) (err error) {
	o := orm.NewOrm()
	v := KpiTask{Id: task.Id}
	if err = o.Read(&v); err == nil {
		_, err = o.Update(task, field...)
		if err != nil {
			beego.Error(err.Error())
			return
		}
	} else {
		beego.Error(err.Error())
	}
	return
}
