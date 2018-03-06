package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/kataras/go-errors"
	"time"
)

type UserFace struct {
	Id         int64  `json:"id,omitempty" orm:"column(id);auto"`
	Uid        int    `json:"uid" orm:"column(uid)"`
	Status     int    `json:"status" orm:"column(status)"`
	Path       string `json:"path" orm:"column(path)"`
	Remarks    string `json:"remarks" orm:"column(remarks)"`
	Name       string `json:"name" orm:"column(name)"`
	CreateTime int64  `json:"create_time" orm:"column(create_time)"`
	UpdateTime int64  `json:"update_time" orm:"column(update_time)"`
}

type FaceJson struct {
	Uid  int    `json:"uid"`
	Img  string `json:"img"`
	Tag  int    `json:"tag"`
	Code string `json:"code"`
}

func (t *UserFace) TableName() string {
	return "user_face"
}

func init() {
	orm.RegisterModel(new(UserFace))
}

//判断用户是否认证过
//3 未查到该用户绑定信息 2 正在审核  1审核通过  3未注册 4不通过
func GetExistFaceByUid(id int) (int) {
	o := orm.NewOrm()
	user := &UserFace{Uid: id}
	if err := o.Read(user, "uid"); err != nil {
		return 3
	}
	return user.Status
}

//通过uid查询
func GetFacePathByUid(id int) (err error, path string) {
	o := orm.NewOrm()
	user := &UserFace{Uid: id}
	if err = o.Read(user, "uid"); err != nil {
		return
	}
	path = user.Path
	return
}

//添加用户
func AddUserFace(u *UserFace) (err error) {
	o := orm.NewOrm()
	user := &UserFace{Uid: u.Uid}
	var name string
	err, name = GetNickNameById(u.Uid)
	if err != nil {
		return
	}
	if err = o.Read(user, "uid"); err == nil {
		if user.Status == 3 {
			user.CreateTime = time.Now().Unix()
			user.Path = u.Path
			user.Status = 2
			user.Remarks = ""
			var index int64
			index, err = o.Update(user, "update_time", "path", "status", "remarks")
			if index <= 0 || err != nil {
				err = errors.New("插入失败")
				return
			}
			return
		}
		err = errors.New("该用户已经存在人脸认证信息")
		return
	}
	u.CreateTime = time.Now().Unix()
	var index int64
	u.Name = name
	index, err = o.Insert(u)
	if index <= 0 || err != nil {
		err = errors.New("插入失败")
		return
	}

	return
}

//修改认证状态
func ChangeStatusByUid(uid int, status int, remarks string) (err error) {
	o := orm.NewOrm()
	user := &UserFace{Uid: uid}
	if err = o.Read(user, "uid"); err != nil {
		err = errors.New("该用户没有信息")
		return
	}
	//存在
	user.UpdateTime = time.Now().Unix()
	user.Status = status
	user.Remarks = remarks
	var index int64
	index, err = o.Update(user)
	if index <= 0 || err != nil {
		err = errors.New("修改失败")
		return
	}
	return
}
