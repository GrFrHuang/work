package models

import (
	"testing"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"encoding/json"
	"strconv"
	"time"
)
type Tmpcontract struct {
	Id           int     `json:"id,omitempty" orm:"column(id);auto"`
	Ladders       string     `json:"ladders,omitempty" orm:"column(ladders);null"`
	GameName       string     `json:"game_name,omitempty" orm:"column(game_name);null"`
	Cname       string     `json:"cname,omitempty" orm:"column(cname);null"`
	Tname       string     `json:"tname,omitempty" orm:"column(tname);null"`
	Ratiostr	string `json:"ratiostr,omitempty" orm:"column(ratiostr);null"`

}
type TmpRatio struct {
	Ratio       float32     `json:"ratio,omitempty" orm:"column(ratio);null"`
	SlottingFee       float32     `json:"slotting_fee,omitempty" orm:"column(slotting_fee);null"`

}

func (t *Tmpcontract) TableName() string {
	return "tmpcontract"
}

func init() {
	orm.RegisterModel(new(Tmpcontract))
}

func TestLadders(t *testing.T) {


	contract := []Tmpcontract{}
	_, err := orm.NewOrm().Raw("SELECT a.ladders,b.`game_name`,c.`name` as cname,d.`name` as tname FROM contract a LEFT JOIN game b ON " +
		"a.`game_id`=b.`game_id` LEFT JOIN channel c ON a.channel_code = c.cp  LEFT JOIN `types` d ON a.`state`=d.`id` WHERE " +
		"a.company_type=1 AND a.game_id IN " +
		"(816,828,835,840,848,862,882,888,914,915,918,923,928,932,944,967,977,982,1009,1020,1046,1055,1060,1062,1063,1065,1073," +
		"1075,1078,1079,1080,1083,1096,1099,1104)",
		).QueryRows(&contract)
	if err != nil {
		return
	}
	for k,v := range(contract) {
		if v.Ladders  != "" {
			f := []TmpRatio{}
			json.Unmarshal([]byte(v.Ladders), &f)
			s,_ := json.Marshal(f)
			contract[k].Ratiostr = string(s)
			v.Ratiostr = string(s)
			beego.Warning(f)

			tmpc := new(Tmpcontract)
			tmpc.Ladders = v.Ladders
			tmpc.GameName = v.GameName
			tmpc.Cname = v.Cname
			tmpc.Tname = v.Tname
			tmpc.Ratiostr = v.Ratiostr
			orm.NewOrm().Insert(tmpc)
		}
	}
	beego.Warning(contract)
}
