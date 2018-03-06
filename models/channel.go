package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/utils"
	"reflect"
	"strings"
	"fmt"
)

type Channel struct {
	Id        int    `json:"id,omitempty" orm:"column(id);null" `
	ChannelId int    `json:"channel_id,omitempty" orm:"column(channel_id);null;size(100)" form:"channel_id"  valid:"Required"` //渠道ID
	Name      string `json:"name,omitempty" orm:"column(name);null;size(100)" form:"Name"  valid:"Required"`
	Cp        string `json:"cp,omitempty" orm:"column(cp);null;size(255)" form:"cp"  valid:"Required"`                    //渠道标识
	Platform  string `json:"platform,omitempty" orm:"column(platform);null;size(255)" form:"platform"  valid:"Required" ` //平台
}

func (t *Channel) TableName() string {
	return "channel"
}

func init() {
	orm.RegisterModel(new(Channel))
}

//获取所有合作中的渠道
func GetAllChannels() (channels []Channel, total int64, err error) {
	o := orm.NewOrm()
	total, err = o.Raw(" SELECT channel_id,name,cp FROM channel a LEFT JOIN channel_company b ON a.cp = b.channel_code WHERE b.cooperate_state=1").QueryRows(&channels)
	return
}

func GetAllChannel(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Channel))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, v == "true" || v == "1")
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

	var l []Channel
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

// 根据渠道商公司主体id，获取该主体的渠道以及渠道负责人
func GetChannelAndPeople(companyId int64) (channelNameWithPepople string) {
	var channelCompanys []orm.Params

	sql := "SELECT a.channel_code,b.name as ChannelName,c.nickname as YunduanResPerName ,d.nickname as YouliangResPerName " +
		"FROM channel_company a " +
		"LEFT JOIN channel b ON a.channel_code=b.cp " +
		"LEFT JOIN `user` c ON a.yunduan_responsible_person=c.id " +
		"LEFT JOIN `user` d ON a.youliang_responsible_person=d.id " +
		"WHERE company_id = ? "
	orm.NewOrm().Raw(sql, companyId).Values(&channelCompanys)

	var names []string
	for _, company := range channelCompanys {
		yunName := company["YunduanResPerName"]
		if yunName == nil {
			yunName = "无"
		}
		youName := company["YouliangResPerName"]
		if youName == nil {
			youName = "无"
		}
		//<span style='color: red'>%s~%s</span>
		name := fmt.Sprintf("%v（【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v)", company["ChannelName"], yunName, youName)
		fmt.Printf("name:%v\n", name)
		names = append(names, name)
	}

	return strings.Join(names, "；")
}

// 根据渠道code,获取该渠道的商务负责人
func GetPeopleByChannelCode(channelCode string) (people string) {
	var channelCompanys []orm.Params

	sql := "SELECT a.channel_code,c.nickname as YunduanResPerName ,d.nickname as YouliangResPerName " +
		"FROM channel_company a " +
		"LEFT JOIN `user` c ON a.yunduan_responsible_person=c.id " +
		"LEFT JOIN `user` d ON a.youliang_responsible_person=d.id " +
		"WHERE a.channel_code = ? "
	orm.NewOrm().Raw(sql, channelCode).Values(&channelCompanys)

	var names []string
	for _, company := range channelCompanys {
		yunName := company["YunduanResPerName"]
		if yunName == nil {
			yunName = "无"
		}
		youName := company["YouliangResPerName"]
		if youName == nil {
			youName = "无"
		}
		name := fmt.Sprintf("【<span style='color: green'>云端</span>】%v【<span style='color: FF9900'>有量</span>】%v", yunName, youName)
		fmt.Printf("name:%v\n", name)
		names = append(names, name)
	}

	return strings.Join(names, "；")
}

//根据渠道id获取渠道名
func GetChannelNameById(channelid int) (string) {
	o := orm.NewOrm()
	channel := Channel{ChannelId: channelid}
	o.Read(&channel)
	return channel.Name
}

//根据渠道cp获取渠道名
func GetChannelNameByCp(cp string) (string, error) {
	o := orm.NewOrm()
	channel := Channel{}
	qs := o.QueryTable(new(Channel))
	err := qs.Filter("CP__exact", cp).One(&channel)
	//err := o.Read(&channel, "Cp")
	if err != nil {
		return "", err
	}
	return channel.Name, nil
}

func GetChannelIdByCp(cp string) (int, error) {
	o := orm.NewOrm()
	channel := Channel{Cp: cp}
	err := o.Read(&channel, "Cp")
	if err != nil {
		return 0, err
	}
	return channel.ChannelId, nil
}

func GetChannelByChannelId(channelId int) (channel Channel, err error) {
	o := orm.NewOrm()
	channel.ChannelId = channelId
	err = o.Read(&channel, "Channelid")
	if err != nil {
		return
	}
	return
}

//根据游戏id，获取该游戏可以下发的渠道（排除该游戏已发渠道）
func GetAddChannels(gameId int) (channels []Channel, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("SELECT a.channel_id,a.name,a.cp FROM channel a LEFT JOIN channel_company b ON a.cp = b.channel_code "+
		"WHERE a.cp NOT IN (SELECT channel_code FROM channel_access WHERE game_id = ?) AND b.cooperate_state=1 ",
		gameId).QueryRows(&channels)
	return
}

func GetChannelByChannelCode(cp string) (channel *Channel, err error) {
	channel = &Channel{}
	o := orm.NewOrm()
	channel.Cp = cp
	err = o.Read(channel, "cp")
	if err != nil {
		return
	}
	return
}

func GetChannelByChannelName(name string) (channel *Channel, err error) {
	channel = &Channel{}
	o := orm.NewOrm()
	channel.Name = name
	err = o.Read(channel, "name")
	if err != nil {
		return
	}
	return
}

func GetChannelNameByChannelCode(code string) (name string, err error) {
	table := "channelCode2Name"
	key := code
	name, err = utils.Redis.HMGETOne(table, key)
	if err != nil {
		return
	}

	if name != "" {
		return
	}

	var channel []Channel
	_, err = orm.NewOrm().QueryTable("channel").All(&channel, "cp", "name")
	if err != nil {
		return
	}

	data := make(map[string]interface{}, len(channel))
	for _, v := range channel {
		data[v.Cp] = v.Name
	}

	name, _ = util.Interface2String(data[key], false)
	if name == "" {
		err = errors.New("404")
		return
	}

	err = utils.Redis.HMSETALL(table, data, 2*60)
	if err != nil {
		return
	}

	return
}
