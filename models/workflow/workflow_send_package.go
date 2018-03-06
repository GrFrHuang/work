package workflow

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type WorkflowSendPackage struct {
	Id              int    `orm:"column(id);auto"`
	ChannelInfoId   int    `orm:"column(channel_info_id);null" description:"渠道信息ID"`
	ContractInfoId  int    `orm:"column(contract_info_id);null" description:"合同信息ID"`
	UserId          string `orm:"column(user_id);size(100);null" description:"用户操作（1,2,3）顺序执行"`
	TaskId          int    `orm:"column(task_id);null" description:"任务ID"`
	CreateTime      int64  `orm:"column(create_time);null" description:"创建时间"`
	UpdateTime      int64  `orm:"column(update_time);null" description:"创建时间"`
	CurrentUserIds  string `orm:"column(current_user_ids);null" description:"当前操作的用户"`
	Current_user_id string `orm:"column(current_user_id);null description:"当前操作的UID""`
	Status 			int    `orm:"-" json:"status"`
	Remarks			string `orm:"-" json:"remarks"`
}

func (t *WorkflowSendPackage) TableName() string {
	return "workflow_send_package"
}

func init() {
	orm.RegisterModel(new(WorkflowSendPackage))
}

func AddWorkflowSendPackage(m *WorkflowSendPackage) (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(m)
	return
}

func ChangeWorkflowInfoById(m *WorkflowSendPackage) (err error) {
	o := orm.NewOrm()
	var i int64
	m.UpdateTime = time.Now().Unix()
	i, err = o.Update(m)
	if i <= 0 || err != nil {
		err = errors.New("修改发包工作流信息失败")
		return
	}
	return
}

func GetWorkflowSendPackageByTaskId(tid int) (v *WorkflowSendPackage, err error) {
	o := orm.NewOrm()
	v = &WorkflowSendPackage{TaskId: tid}
	if err = o.Read(v, "task_id"); err == nil {
		return v, nil
	}
	return nil, errors.New("未找到该条任务")
}

func GetCurretUserIdsById(send WorkflowSendPackage) (WorkflowSendPackage, error) {
	o := orm.NewOrm()
	err := o.Read(&send)
	return send, err

}

func IsHandle(sendid int, userid int) (bool, error) {
	send := WorkflowSendPackage{
		Id: sendid,
	}
	send, err := GetCurretUserIdsById(send)
	if err != nil {
		return false, err
	}
	ids := strings.Split(send.CurrentUserIds, ",")
	for _, v := range ids {

		if v == strconv.Itoa(userid) {
			return true, nil
		}
	}
	return false, nil
}
