package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type KpiChildTask struct {
	Id           int     `orm:"column(id);auto" json:"id,omitempty"`
	TaskName     string  `orm:"column(task_name);size(50);null" description:"任务名" json:"task_name,omitempty"`
	Period       string  `orm:"column(period);size(20);null" description:"任务时长" json:"period,omitempty"`
	ProgressRate float64 `orm:"column(progress_rate);null;digits(4);decimals(2)" description:"进度" json:"progress_rate,omitempty"`
	Remark       string  `orm:"column(remark);size(255);null" description:"备注" json:"remark,omitempty"`
	Score        float64 `orm:"column(score);null;digits(4);decimals(2)" description:"分数" json:"score,omitempty"`
	Annotations  string  `orm:"column(annotations);size(255);null" description:"领导批注" json:"annotations,omitempty"`
	CreateTime   int     `orm:"column(create_time);null" description:"创建时间" json:"create_time,omitempty"`
	UpdateTime   int     `orm:"column(update_time);null" description:"修改时间" json:"update_time,omitempty"`
	TaskId       int     `orm:"column(task_id);null" description:"修改时间" json:"task_id,omitempty"`
	Flag         int8    `orm:"column(flag);null" description:"0.被动任务  1.主动任务" json:"flag,omitempty"`
}

type ChildTasks struct {
	PassiveTasks    []*KpiChildTask `json:"passive_tasks,omitempty"`    //被动考核
	InitiativeTasks []*KpiChildTask `json:"initiative_tasks,omitempty"` //主动考核
}

func (t *KpiChildTask) TableName() string {
	return "kpi_child_task"
}

func init() {
	orm.RegisterModel(new(KpiChildTask))
}

// AddKpiChildTask insert  multiple KpiChildTasks into database and returns
func AddKpiChildTasks(m []*KpiChildTask) (err error) {
	o := orm.NewOrm()
	for k := range m {
		m[k].CreateTime = int(time.Now().Unix())
	}
	_, err = o.InsertMulti(len(m), m)
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

// GetKpiChildTaskById retrieves KpiChildTask by Id. Returns error if
// Id doesn't exist
func GetKpiChildTaskById(id int) (v *KpiChildTask, err error) {
	o := orm.NewOrm()
	v = &KpiChildTask{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateKpiChildTask updates KpiChildTask by Id and returns error if
// the record to be updated doesn't exist
func UpdateKpiChildTaskById(m *KpiChildTask) (err error) {
	o := orm.NewOrm()
	v := KpiChildTask{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteKpiChildTask deletes KpiChildTask by Id and returns error if
// the record to be deleted doesn't exist
func DeleteKpiChildTask(id int) (err error) {
	o := orm.NewOrm()
	v := KpiChildTask{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&KpiChildTask{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Get child task by parent task id
func AddChildTask(v *KpiTask) (err error) {
	o := orm.NewOrm()
	kct := []*KpiChildTask{}
	_, err = o.QueryTable("kpi_child_task").Filter("task_id", v.Id).All(&kct, "task_name")
	if err != nil {
		beego.Error(err)
		return
	}
	v.KpiChildTasks = kct
	return
}

func UpdateChildTasks(task *KpiTask) (err error) {
	o := orm.NewOrm()
	ids := make(map[interface{}]bool)
	list := orm.ParamsList{}
	total, err := o.QueryTable("kpi_child_task").Filter("task_id", task.Id).ValuesFlat(&list, "id")
	if err != nil && err != orm.ErrNoRows {
		beego.Error(err)
		return
	}
	for _, v := range list {
		ids[v] = true
	}

	if int(total) == len(task.KpiChildTasks) {
		//Verify whether all child task id are exist or not
		for _, v := range task.KpiChildTasks {
			if _, ok := ids[v.Id]; !ok {
				break
			}
		}
		return
	}
	o.Begin()
	//Delete all then add
	_, err = o.QueryTable("kpi_child_task").Filter("task_id", task.Id).Delete()
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return
	}
	_, err = o.InsertMulti(len(task.KpiChildTasks), task.KpiChildTasks)
	if err != nil {
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

//Examine task
func ExamineChildTask(tasks []*KpiChildTask) (err error) {
	o := orm.NewOrm()
	for _, v := range tasks {
		v.UpdateTime = int(time.Now().Unix())
		_, err = o.Update(v, "score", "annotations", "update_time")
		if err != nil {
			beego.Error(err)
			return
		}
	}
	t, err := GetKpiTaskById(tasks[0].TaskId)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	t.AuditState = 1
	t.UpdateTime = int(time.Now().Unix())
	err = UpdateKpiTaskFields(t, "audit_state", "update_time")
	if err != nil {
		beego.Error(err.Error())
	}
	return
}

func WriteTask(childTasks *ChildTasks, id int) (err error) {
	o := orm.NewOrm()
	for _, v := range childTasks.PassiveTasks {
		v.UpdateTime = int(time.Now().Unix())
		_, err = o.Update(v, "period", "progress_rate", "remark", "update_time")
		if err != nil {
			beego.Error(err)
			return
		}
	}
	//Judge whether add initiative tasks or not
	if len(childTasks.InitiativeTasks) > 0 {
		for _, v := range childTasks.InitiativeTasks {
			v.Flag = 1
			v.TaskId = id
			task := KpiChildTask{
				Id: v.Id,
			}
			if err = o.Read(&task); err == nil {
				v.UpdateTime = int(time.Now().Unix())
				_, err = o.Update(v, "period", "progress_rate", "remark", "update_time")
				if err != nil {
					beego.Error(err)
					return
				}
			}
			if err == orm.ErrNoRows {
				v.CreateTime = int(time.Now().Unix())
				_, err = o.Insert(v)
				if err != nil {
					beego.Error(err)
					return
				}
			}
			if err != orm.ErrNoRows && err != nil {
				beego.Error(err)
			}
		}
	}
	return
}

func AddTotalScore(kts []*KpiTask) ([]*KpiTask, error) {
	var ids []int
	var kct []KpiChildTask
	var scores = make(map[int]float64)
	var str string
	for _, v := range kts {
		ids = append(ids, v.Id)
	}
	for k := range kts {
		if k == len(kts)-1 {
			str = str + "?"
			break
		}
		str = "?," + str
	}
	o := orm.NewOrm()

	_, err := o.Raw("select task_id, sum(score) as score from kpi_child_task where task_id in ("+str+") group by task_id ;", ids).QueryRows(&kct)
	if err != nil {
		beego.Error(err)
		return nil, nil
	}
	for _, v := range kct {
		scores[v.TaskId] = v.Score
	}
	for k, v := range kts {
		kts[k].TotalScore = scores[v.Id]
	}
	return kts, nil
}

//Get Child kpi Task by parent task id
func GetChildTaskByTaskId(taskId int) (childTasks []*KpiChildTask, err error) {
	return
}
