package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
	"time"
	"strconv"
	"kuaifa.com/kuaifa/work-together/uc_client"
)

type UserFaceVerify struct {
	Id         int    `orm:"column(id);auto" json:"id,omitempty"`
	Uid        int    `orm:"column(uid);null" description:"用户ID" json:"uid,omitempty"`
	Tag        int    `orm:"column(tag);null" description:"1打款 2登录" json:"tag,omitempty"`
	Name       string `orm:"column(name);size(50);null" description:"用户名" json:"name,omitempty"`
	Image      string `orm:"column(image);size(100);null" description:"人脸记录" json:"image,omitempty"`
	Remarks    string `orm:"column(remarks);size(255);null" description:"备注" json:"remarks,omitempty"`
	CreateTime int64  `orm:"column(create_time);null" description:"创建时间" json:"create_time,omitempty"`

	Token string `json:"token,omitempty" orm:"-"`
	Code  string `json:"code,omitempty" orm:"-"`
}

func (t *UserFaceVerify) TableName() string {
	return "user_face_verify"
}

func init() {
	orm.RegisterModel(new(UserFaceVerify))
}

func AddUserCodeLogin(m *UserFaceVerify) (error) {
	if m.Tag != 2 {
		return errors.New("tag !=2")
	}
	vcode, err := utils.Redis.GET(m.Code)
	if err != nil {
		return err
	}
	if len(vcode) <= 0 {
		return errors.New("二维码已过期")
	}
	utils.Redis.DEL(m.Code)
	u, err := GetUserInfoByToken(m.Token)
	if err != nil {
		return errors.New("token invalid")

	}
	res, err := uc_client.CheckAccessToken(m.Token, u.Name)
	if err != nil && !res {
		return errors.New("token invalid")
	}

	user := fmt.Sprintf(`{"email": "%s","password": "%s"}`, u.Name, m.Token)
	fmt.Println(user)
	if err = utils.Redis.SET(m.Code, user, 2*60); err != nil {
		return err
	}
	o := orm.NewOrm()
	m.Uid = u.Id
	m.Name = u.Nickname
	m.CreateTime = time.Now().Unix()
	_, err = o.Insert(m)
	return err
}

// AddUserFaceVerify insert a new UserFaceVerify into database and returns
// last inserted Id on success.
func AddUserFaceVerify(m *UserFaceVerify) (error) {
	if m.Tag != 3 {
		return errors.New("tag !=3")
	}
	vcode, err := utils.Redis.GET("face_check")
	if err != nil {
		return err
	}
	if vcode != m.Code {
		return errors.New("二维码已过期或不正确")
	}
	o := orm.NewOrm()
	as, err := NewAsset().SaveAssetBase64Img(m.Image, m.Uid)
	if err != nil {
		return err
	}
	m.Image = strconv.FormatInt(as.Id, 10)
	err, name := GetNickNameById(m.Uid)
	if err != nil {
		return err
	}
	m.Name = name
	m.CreateTime = time.Now().Unix()
	_, err = o.Insert(m)
	return err
}
