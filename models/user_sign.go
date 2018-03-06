package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"kuaifa.com/kuaifa/work-together/utils"
)

type UserSign struct {
	Id       int    `orm:"column(id);auto" json:"id"`
	Uid      int    `orm:"column(uid);null" description:"用户ID" json:"uid"`
	Name     string `orm:"column(name);size(50);null" json:"name"`
	Status   int    `orm:"column(status);null" description:"用户签到审核" json:"status"`
	Path     string `orm:"column(path);size(100);null" json:"path"`
	Tag      int    `orm:"column(tag);null" json:"tag"`
	Remarks  string `json:"remarks" orm:"column(remarks)"`
	Date     string `orm:"column(date);size(50);null" json:"date"`
	SignTime string `orm:"column(sign_time);size(100);null" description:"签到时间" json:"sign_time"`
}


func (t *UserSign) TableName() string {
	return "user_sign"
}

func init() {
	orm.RegisterModel(new(UserSign))
}

// AddUserSign insert a new UserSign into database and returns
// last inserted Id on success.
func AddUserSign(uid int, tag int, img string, code string) (error) {
	o := orm.NewOrm()
	if tag == 1 {
		vcode, err := utils.Redis.GET("morning_sign")
		if err != nil {
			return err
		}
		if vcode != code {
			return errors.New("二维码已过期或不正确")
		}
	}
	//验证是否已经签到
	num, err := o.QueryTable("user_sign").Filter("uid", uid).Filter("tag", tag).Filter("sign_time__icontains", time.Now().Format("2006-01-02")).Count()
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("该用户已经签到")
	}
	if err := utils.FaceVerify(strconv.Itoa(uid), img); err != nil {
		return errors.New("授权服务器:" + err.Error())
	}
	var sign UserSign
	sign.Uid = uid
	sign.Tag = tag
	sign.Status = 2
	sign.SignTime = time.Now().Format("2006-01-02 15:04:05")
	sign.Date = time.Now().Format("2006-01-02")
	err, name := GetNickNameById(uid)
	if err != nil {
		return err
	}
	sign.Name = name
	//添加图片地址
	as, err := NewAsset().SaveAssetBase64Img(img, uid)
	if err != nil {
		return err
	}
	sign.Path = strconv.FormatInt(as.Id, 10)
	index, err := o.Insert(&sign)
	if index <= 0 || err != nil {
		return err
	}
	return nil
}

//修改认证状态
func ChangeSignStatusByUid(uid int, status int, remarks string) (err error) {
	o := orm.NewOrm()
	user := &UserSign{Uid: uid}
	if err = o.Read(user, "uid"); err != nil {
		err = errors.New("该用户没有信息")
		return
	}
	//存在
	user.Status = status
	var index int64
	user.Remarks = remarks
	index, err = o.Update(user)
	if index <= 0 || err != nil {
		err = errors.New("修改失败")
		return
	}
	return
}
