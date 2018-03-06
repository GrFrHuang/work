package workflow

import (
	"sync"
	"strings"
	"strconv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"errors"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
	"kuaifa.com/kuaifa/work-together/tool"
)

type WorkflowTask struct {
	Id                  int    `orm:"column(id);auto" json:"id,omitempty"`
	TaskName            string `orm:"column(task_name);size(100);null" description:"任务名称" json:"task_name,omitempty"`
	WfNameId            int    `orm:"column(wf_name_id);null" description:"工作流ID" json:"wf_name_id,omitempty"`
	CurrentProgress     string `orm:"column(current_progress);null" description:"当前的进度" json:"current_progress,omitempty"`
	CurrentProgressName string `orm:"column(current_progress_name);null" description:"当前的进度" json:"current_progress_name,omitempty"`
	CurrentSuccess      int    `orm:"column(current_success);null" description:"当前完成的" json:"current_success,omitempty"`
	CreateTime          int64  `orm:"column(create_time);null" description:"创建时间" json:"create_time,omitempty"`
	UpdateTime          int64  `orm:"column(update_time);null" description:"修改时间" json:"update_time,omitempty"`
	Status              int    `orm:"column(status);null" description:"状态（1成功 2失败 3挂起 4创建 5其他,6 在流程中）" json:"status,omitempty"`
	Remarks             string `orm:"column(remarks);size(100);null" description:"备注信息" json:"remarks,omitempty"`
	CreateName          string `orm:"column(create_name);size(100);null" description:"创建人名字" json:"create_name,omitempty"`
	UpdateName          string `orm:"column(update_name);size(100);null" description:"更新人名字" json:"update_name,omitempty"`
	GameName            string `orm:"column(game_name);size(100);null" description:"游戏名称" json:"game_name,omitempty"`
	ChannelName         string `orm:"column(channel_name);size(100);null" description:"渠道名称" json:"channel_name,omitempty"`
	Step                int    `orm:"column(step);null" json:"step"`

	UserId       int    `orm:"-" json:"user_id,omitempty"`
	GameId       int    `orm:"-" json:"game_id,omitempty"`
	ChannelCode  string `orm:"-" json:"channel_code,omitempty"`
	DepartmentId int    `orm:"-" json:"department_id,omitempty"`
}

type Users struct {
	Uid          int    `json:"uid"`
	DepartmentId int    `json:"department_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}

type ChangeProgressInput struct {
	User         []Users `json:"user"`
	Remarks      string  `json:"remarks"`
	PerateUid    int     `json:"perate_uid"`
	TaskId       int     `json:"task_id"`
	DepartmentId string  `json:"department_id"`
	NextNodeId   string  `json:"next_node_id"`
	Step         int     `json:"step"`
	Rollback     int     `json:"rollback"`
	UserId       int     `json:"user_id"`
}

type ConditionTask struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

var sy sync.Mutex

func (t *WorkflowTask) TableName() string {
	return "workflow_task"
}

func init() {
	orm.RegisterModel(new(WorkflowTask))
}

// AddWorkflowTask insert a new WorkflowTask into database and returns
// last inserted Id on success.
func AddWorkflowTask(m *WorkflowTask) (error, int) {
	o := orm.NewOrm()
	task := WorkflowTask{TaskName: m.TaskName}
	if err := o.Read(&task, "task_name"); err == nil {
		err = errors.New("该条任务已经存在")
		return err, 0
	}
	err := o.Begin()
	if err != nil {
		return err, 0
	}
	sy.Lock()
	defer sy.Unlock()
	//通过部门ID获取结点信息
	info, err := GetWorkflowNodeByDepartmentId(m.DepartmentId, m.Step)
	if err != nil {
		o.Rollback()
		return err, 0
	}
	//判断是否有权限
	if len(info.CurrentNodeId) != 0 {
		err = errors.New("没有权限创建工作流")
		o.Rollback()
		return err, 0

	}
	err, name := models.GetNickNameById(m.UserId)
	if err != nil {
		o.Rollback()
		return err, 0
	}
	names := strings.Split(m.TaskName, "&")
	if len(names) != 2 {
		o.Rollback()
		return errors.New("未获取到渠道或者游戏名字"), 0
	}
	//插入任务表
	m.CreateTime = time.Now().Unix()
	m.Status = 4
	m.ChannelName = names[0]
	m.GameName = names[1]
	m.CreateName = name
	m.CurrentProgress = strconv.Itoa(info.Id)
	m.CurrentProgressName = info.NodeName
	id, err := o.Insert(m)
	if err != nil || id == 0 {
		err = errors.New("添加失败" + err.Error())
		o.Rollback()
		return err, 0
	}
	//锁定渠道信息
	var ChannelInfoId int
	var ContractInfoId int
	ChannelInfoId, err = models.ChangeWorkflowStatusByid(m.GameId, 1, m.ChannelCode)
	if err != nil {
		o.Rollback()
		return err, 0
	}
	//获取合同ID
	ContractInfoId, err = models.GetContractByGameid(m.GameId, m.ChannelCode)
	if err != nil {
		o.Rollback()
		return err, 0
	}

	//添加到发包流程表
	var send WorkflowSendPackage
	send.ChannelInfoId = ChannelInfoId
	send.ContractInfoId = ContractInfoId
	send.CreateTime = time.Now().Unix()
	send.CurrentUserIds = strconv.Itoa(m.UserId)
	send.UserId = strconv.Itoa(m.UserId)
	send.TaskId = int(id)
	err = AddWorkflowSendPackage(&send)
	if err != nil {
		o.Rollback()
		err = errors.New("添加失败" + err.Error())
		return err, 0
	}
	err = AddLog(m.UserId, send.TaskId, send.CreateTime, send.Id, task.TaskName, "", 1)
	if err != nil {
		o.Rollback()
		err = errors.New("日志添加失败" + err.Error())
		return err, 0
	}
	err = o.Commit()
	if err != nil {
		return err, 0
	}
	return nil, int(id)
}

//通过条件查询数据
func GetAllByCondition(conditions []string) (interface{}, error) {
	var data = map[string][]orm.Params{}
	for _, v := range conditions {
		o := orm.NewOrm()
		var lists []orm.Params
		sql := fmt.Sprintf("SELECT id,%s  FROM workflow_task GROUP BY %s", v, v)
		num, err := o.Raw(sql).Values(&lists)
		if err == nil && num > 0 {
			data[v] = lists
		}
	}
	return data, nil
}

//修改任务状态
func ChangeStatucbyId(id, status, channeld, contractId int, remarks string) error {
	if status == 1 || status == 2 || status == 3 {
		if err := models.ChangeWorkflowStatus(channeld, 2); err != nil {
			return err
		}
		if status == 2 {
			//删除合同
			if err := models.DelContractByid(contractId); err != nil {
				return err
			}
		}
	}
	var w WorkflowTask
	w.Id = id
	w.UpdateTime = time.Now().Unix()
	w.Status = status
	w.Remarks = remarks
	o := orm.NewOrm()
	i, err := o.Update(&w, "status", "update_time", "remarks")
	if i <= 0 || err != nil {
		return errors.New("修改进度状态失败")
	}

	return nil
}

//修改任务状态
func ChangeCurrentProgressById(id, step int, departmentId, progress, remarks string) error {
	var w WorkflowTask
	w.Id = id
	o := orm.NewOrm()
	if err := o.Read(&w); err != nil {
		return err
	}
	if err := GetWorlflowBodeDepartmentIdsById(w.CurrentProgress, departmentId); err != nil {
		return err
	}
	taskName, err := GetWorkflowNodeByIds(progress)
	if err != nil {
		return err
	}
	w.UpdateTime = time.Now().Unix()
	w.CurrentProgress = progress
	w.Remarks = remarks
	w.CurrentProgressName = taskName
	w.Status = 6
	w.CurrentSuccess = 0
	w.Step = step

	i, err :=
		o.Update(&w, "current_progress", "update_time", "remarks", "current_success", "current_progress_name", "status", "step")
	if i <= 0 || err != nil {
		return errors.New("修改进度信息失败")
	}
	return nil
}

//修改最后更新人信息
func ChangeUpdateNameById(id int, name string) error {
	o := orm.NewOrm()
	var w WorkflowTask
	w.Id = id
	w.UpdateTime = time.Now().Unix()
	w.UpdateName = name
	i, err := o.Update(&w, "update_name", "update_time")
	if i <= 0 || err != nil {
		return errors.New("修改进度信息失败")
	}
	return nil
}

func AddCurrentSuccessById(id int) error {
	var w WorkflowTask
	w.Id = id
	o := orm.NewOrm()
	if err := o.Read(&w); err != nil {
		return err
	}
	w.CurrentSuccess += 1
	w.UpdateTime = time.Now().Unix()
	i, err := o.Update(&w, "current_success", "update_time")
	if err != nil || i <= 0 {
		return errors.New("修改进度完成数信息失败")
	}
	return nil
}

func ChangeUpdateMan(v WorkflowSendPackage) error {
	o := orm.NewOrm()
	work := WorkflowTask{Id: v.TaskId}

	ids := strings.Split(v.CurrentUserIds, ",")
	if len(ids) > 0 {

		//userids := v.UserId
		userids := strings.Split(v.UserId, ",")
		userId, err := strconv.Atoi(userids[len(userids)-1])
		//userId, err := strconv.Atoi(ids[0])
		if err != nil {
			return errors.New("user ID 无效")
		}

		err, userName := models.GetNickNameById(userId)
		if err != nil {
			return errors.New("没有找到该用户")
		}
		err = o.Read(&work, "id")
		if err != nil {
			return errors.New("没有找到该任务信息")
		}
		work.UpdateName = userName
		work.UpdateTime = time.Now().Unix()

		_, err = o.Update(&work)
		if err != nil {
			return errors.New("更新失败")
		}

		next_user_id, _ := strconv.Atoi(ids[0])
		err, next_user := models.GetNickNameById(next_user_id)
		if err != nil {
			return errors.New("下一节点用户未找到")
		}
		err = AddLog(userId, v.TaskId, time.Now().Unix(), v.Id, v.Remarks, next_user, 3)
		if err != nil {
			return err
		}
	}

	return nil

}

//修改当前任务进度
func ChangeProgress(v *ChangeProgressInput) error {
	if len(v.User) == 0 {
		return errors.New("请选择需要转交的用户")
	}
	if v.NextNodeId == "" {
		return errors.New("未上传下一结点信息")
	}
	if v.Step == 0 {
		return errors.New("未获取到步骤信息")
	}
	//添加用户信息到工作流程表
	task, err := GetWorkflowTaskById(v.TaskId)
	if err != nil {
		return err
	}
	if task.Status == 1 {
		return errors.New("当前工作已经完成")
	}
	err, userName := models.GetNickNameById(v.UserId)
	if err != nil {
		return errors.New("未查询到用户名称")
	}
	current := strings.Split(task.CurrentProgress, ",")
	if len(current) == 1 {
		//完成可以移交
		step := getNextNodeStep(v.NextNodeId)
		if err := ChangeCurrentProgressById(v.TaskId, step, v.DepartmentId, v.NextNodeId, v.Remarks); err != nil {
			return err
		}
		for _, value := range v.User {
			//发送邮件
			go func(u Users, task string) {
				err := tool.SendWorkflowEmail(u.Name, task, u.Email)
				fmt.Println(err)
			}(value, task.TaskName)
			//添加日志
			err = AddLog(v.PerateUid, task.Id, time.Now().Unix(), task.WfNameId, v.Remarks, value.Name, 2)
			if err != nil {
				err = errors.New("日志添加失败" + err.Error())
				return err
			}
		}
		err = ChangeUpdateNameById(v.TaskId, userName)
		if err != nil {
			err = errors.New("日志添加失败" + err.Error())
			return err
		}

		return nil
	} else {
		if len(current) <= task.CurrentSuccess {
			//完成可以移交
			if err := ChangeCurrentProgressById(v.TaskId, v.Step, v.DepartmentId, v.NextNodeId, v.Remarks); err != nil {
				return err
			}
			//添加日志（完成该流程日志）

			for _, value := range v.User {
				//发送邮件
				go func(u Users, task string) {
					err := tool.SendWorkflowEmail(u.Name, task, u.Email)
					fmt.Println(err)
				}(value, task.TaskName)
				err = AddLog(v.PerateUid, task.Id, time.Now().Unix(), task.WfNameId, v.Remarks, value.Name, 2)
				if err != nil {
					err = errors.New("日志添加失败" + err.Error())
					return err
				}
				err = ChangeUpdateNameById(v.TaskId, value.Name)
				if err != nil {
					err = errors.New("日志添加失败" + err.Error())
					return err
				}

			}
			err = ChangeUpdateNameById(v.TaskId, userName)
			if err != nil {
				err = errors.New("日志添加失败" + err.Error())
				return err
			}
			return nil
		} else {
			//全部办理人未完全办理
			if err := AddCurrentSuccessById(v.TaskId); err != nil {
				return err
			}
			//添加日志（某部门完成了还有其他地方未完成）
			err = AddLog(task.UserId, task.Id, time.Now().Unix(), task.WfNameId, task.Remarks, "", 3)
			if err != nil {
				err = errors.New("日志添加失败" + err.Error())
				return err
			}
			err = ChangeUpdateNameById(v.TaskId, userName)
			if err != nil {
				err = errors.New("日志添加失败" + err.Error())
				return err
			}
			return errors.New("还有办理人未办理不能转下一步,全部办理后自动提交,请不要重复点击")
		}
	}
}

// GetWorkflowTaskById retrieves WorkflowTask by Id. Returns error if
// Id doesn't exist
func GetWorkflowTaskById(id int) (v *WorkflowTask, err error) {
	o := orm.NewOrm()
	v = &WorkflowTask{Id: id}
	err = o.Read(v)
	if err != nil {
		return
	}
	return
}

func GetChannelInfoIdByTaskId(taskId int) (info interface{}, err error) {
	send, errs := GetWorkflowSendPackageByTaskId(taskId)
	if errs != nil {
		err = errs
		return
	}
	info = send
	return
}
