package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"strconv"
	"strings"
	"testing"
	"github.com/astaxie/beego"
)

//func TestMigration(t *testing.T) {
//	if err := MigrationVerifyChannel(2); err != nil {
//		t.Error(err)
//	}
//}

func TestGetGameNameByGameId(t *testing.T) {
	name, err := GetGameNameByGameId(1107)
	if err != nil {
		t.Error(err)
	}
	t.Log(name)
}

func TestGetUserNameByUserId(t *testing.T) {
	name, err := GetUserNameByUserId(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(name)
}

func BenchmarkGetUserNameByUserId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGameNameByGameId(1107)
	}
}

// 修复 渠道已经对账的预对账单被标记为-1的错误
//func TestFixChannelVerifyGames(t *testing.T) {
//	FixVerifyChannelFormOld()
//}

type NeedFixChannelVerify struct {
	code      string
	date      string
	companyId int
}

// 将指定渠道对账单标记为已开票
func TestFixChannelVerifyStatus(t *testing.T) {
	s := `2017年1月		44937	44937
2016年10月		1881wan	1881玩
2016年10月		caohua	51手游社区
2016年12月		51wan	51玩
2016年11月		7k7ksy	7k7k手游
2016年9月		7u	7u
2017年1月		91wan	91wan
2017年1月		9665	9665游戏
2016年12月		hupu	hupu
2016年9月		oppo	oppo
2017年1月		pkpk	pkpk
2016年12月		pps	pps
2016年12月		ninegame	uc9游
2016年9月		vivo	vivo
2016年11月		aipu	爱谱
2016年12月		egameddai	爱游戏短代
2016年11月		icebird	冰鸟
2016年10月		boyakenuo	博雅科诺
2016年10月		xxzhushou	叉叉助手
2016年11月		dianjin	点金
2016年11月		dyoo	点优
2017年1月		feichen	飞辰
2016年10月		fengbao	风暴
2016年11月		guopan	果盘
2017年1月		haima	海马安卓
2016年11月		cga	浩方
2016年10月		24235	嗨皮游戏
2016年9月		huawei	华为
2016年12月		75757	惠游
2016年11月		kingcheer	鲸旗天下
2017年1月		9quyx	九趣游戏
2016年10月		kuangwan	狂玩
2016年11月		-	狂玩（秒创）
2016年12月		lehihi	乐嗨嗨
2016年12月		lepay	乐视
2016年9月		lewanduo	乐玩多平台
2016年11月		lewanwx	乐玩无线
2016年11月		93636	米冠(多帆网络)
2017年1月		muzhiwan	拇指玩
2017年2月		nifenglin	逆风鳞
2016年9月		nubia	努比亚
2017年2月		papayou	啪啪游
2017年2月		qidianyule	奇点娱乐
2016年11月		qiwan	奇顽游戏
2016年9月		ququ	趣趣
2016年10月		49you	世加游戏
2016年9月		9133	手游帮
2017年1月		shouyoudao	手游岛
2016年12月		giant	手游咖啡
2016年12月		07073sy	数游
2017年2月		shunwan	顺玩
2016年12月		sogou	搜狗
2016年9月		tanwan	贪玩
2016年12月		tuulang	螳螂游戏
2016年11月		-	同步推
2016年9月		wandoujia	豌豆荚
2017年1月		weme	微米
2017年1月		wqwan	我去玩
2016年11月		xiaocaohy	小草互娱
2016年12月		xiaopi	小皮
2016年11月		x7sy	小七手游
2017年2月		xiongbing	雄兵
2016年12月		eagle	雄鹰游戏
2016年11月		yiliu	一六游戏
2017年1月		mushroom	怡扬(蘑菇玩)
2016年12月		youku	优酷
2016年10月		-	悠游网
2016年12月		-	游宝
2016年11月		-	鱼丸
2016年11月		zxgame	掌炫科技`

	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	ss := strings.Split(s, " ")
	log.Info("test", strings.Join(ss, "|"))

	var st []string
	for _, i := range ss {
		if i != "" {
			st = append(st, i)
		}
	}

	log.Info("test", strings.Join(st, "|"))
	l := len(st)

	var xs []NeedFixChannelVerify
	for i := 0; i < l; i += 3 {
		date := st[i]
		date = strings.Replace(date, "年", "-", -1)
		date = strings.Replace(date, "月", "", -1)
		if len(strings.Split(date, "-")[1]) == 1 {
			date = strings.Replace(date, "-", "-0", -1)
		}

		xs = append(xs, NeedFixChannelVerify{
			code: st[i + 1],
			date: date,
		})
	}

	fixChannelVerify2(xs)
	return

	// 找渠道商
	for i, v := range xs {
		code := v.code
		if code == "-" {
			continue
		}
		cp := ChannelCompany{}
		orm.NewOrm().QueryTable("channel_company").Filter("channel_code", code).One(&cp)
		if cp.CompanyId == 0 {
			log.Info("test", "err code:", code)
			continue
		}

		xs[i].companyId = cp.CompanyId
	}

	log.Info("test", xs)
	// 刷新对账单
	for _, v := range xs {
		if v.code == "-" {
			continue
		}
		r, _ := orm.NewOrm().Raw("update verify_channel SET status = 30,remit_company_id = ? where channel_code = ? AND date < ? ",
			v.companyId, v.code, v.date).Exec()
		a, _ := r.RowsAffected()
		log.Info("test", a)
	}
}

// 将渠道对账单标记为已开票
func fixChannelVerify2(xs []NeedFixChannelVerify) {
	var channelCodes []string
	for _, v := range xs {
		if v.code != "-" {
			channelCodes = append(channelCodes, "'" + v.code + "'")
		}
	}
	ids := strings.Join(channelCodes, ",")
	// 状态为20 渠道不在已经修复的渠道之内
	sql := "update verify_channel SET status = 30 where channel_code not in (%s) and status = 20 AND id < 1907 "
	sql = fmt.Sprintf(sql, ids)
	r, _ := orm.NewOrm().Raw(sql).Exec()
	a, _ := r.RowsAffected()
	log.Info("test", a)

}

// 获取公司id
//func TestGetCompanyIds(t *testing.T) {
//	s := `北京爱视游科技文化有限公司
//北京安趣科技股份有限公司  安趣
//北京大燕互娱科技有限公司
//北京当乐信息技术有限公司   当乐
//北京多点在线科技有限公司  夜神
//北京广狐信息技术有限公司  狐狸助手
//北京海誉动想科技股份有限公司  海马玩
//北京华艺创梦科技有限公司  找乐助手
//北京机锋科技有限公司  机锋
//北京聚乐网络科技有限公司
//北京乐盟互动科技有限公司
//北京蘑菇互娱科技有限公司  蘑菇玩
//北京拇指玩科技有限公司  拇指玩
//北京奇顽科技有限公司
//北京趣玩互动信息技术有限公司  全民助手
//北京搜狗网络技术有限公司
//北京勇士互娱科技有限公司
//北京掌中悦动科技有限公司  游龙聚会
//长沙超神网络有限公司
//长沙七丽网络科技有限公司
//成都盟宝互动科技有限公司  盟宝（麟游SDK）
//成都魔方在线科技有限公司
//成都闲玩网络科技有限公司
//成都游宝科技有限公司  游宝
//大胜数娱（天津）科技有限公司  海马玩
//福建创意嘉和软件有限公司  松果
//福建点金网络科技有限公司
//福建乐游网络科技有限公司  摩格
//福建天趣网络科技有限公司
//福建游龙网络科技有限公司
//福州靠谱网络有限公司  靠谱
//福州市卓普信息技术有限公司
//广东安久科技有限公司  游戏fan
//广东秒创网络科技股份有限公司  秒创（狂玩）
//广西一六游网络科技有限公司
//广州道先网络科技有限公司  大麦助手
//广州火烈鸟网络科技有限公司   果盘叉叉助手
//广州妙乐网络科技有限公司
//广州沙巴克网络科技有限公司    TT语音
//广州闪趣网络科技有限公司
//广州市幻动网络科技有限责任公司
//广州市久邦数码科技有限公司    3G门户
//广州水煮信息科技有限公司
//广州唐腾信息技术有限公司  新芒果玩
//广州想玩信息科技有限公司  迷思
//广州小朋网络科技有限公司  朋友玩
//广州雄兵网络科技有限公司
//广州耀玩网络科技有限公司  手游狗
//广州裕际网络科技有限公司
//广州云翎网络科技有限公司  乐非凡
//贵州指趣网络科技有限公司
//海南汉风科技有限公司
//杭州宝莲灯科技有限公司
//杭州宝亮网络科技有限公司  宝亮
//杭州飞辰网络科技有限公司
//杭州快定网络股份有限公司  七匣子
//河北梦娱星空网络科技有限公司  7U
//湖南奇葩乐游网络科技有限公司
//江西巨网科技股份有限公司  齐齐乐发发
//龙川奇点网络科技有限公司
//绵阳友游网络科技有限公司
//木蚂蚁（北京）科技有限公司
//南京规贝软件科技有限公司 茄子 江苏猎宝
//南京直立行走网络科技有限公司
//山东柯烁网络科技有限公司
//上海冰穹网络科技有限公司  钱宝
//上海都玩网络科技有限公司   都玩93
//上海范特西网络科技有限公司
//上海弘贯网络科技有限公司   6K玩
//上海晋昶网络科技有限公司  49游还是世加游戏
//上海绿岸网络科技股份有限公司
//上海迈微软件科技有限公司
//上海邑世网络科技有限公司
//上海瑛麒动漫科技有限公司
//上海优扬新媒信息技术有限公司  悠悠村
//上海游戏多网络科技股份有限公司  游戏多
//上海娱嘉网络科技有限公司
//上海云群实业网络科技有限公司
//上海掌悦广告有限公司
//上海掌越科技有限公司
//上饶市奇天乐地科技有限公司
//深圳豪邦网络有限公司  17168豪邦
//深圳乐玩无限科技有限公司
//深圳泡椒思志信息技术有限公司
//深圳尚米网络技术有限公司
//深圳市创想天空科技股份有限公司
//深圳市乐巢科技有限公司
//深圳市拇指游玩科技有限公司
//深圳市全民点游科技有限公司
//深圳市迅雷网络技术有限公司
//深圳市值尚互动科技有限公司
//视娱（天津）网络科技有限公司
//四川九趣网络科技有限公司 1883
//四三九九网络股份有限公司
//苏州蜗牛数字科技股份有限公司
//宿迁楚风互娱科技有限公司
//天津安果科技有限公司
//天行（上海）网络科技有限公司  笨手机
//威海龙必达信息技术有限公司  爱玩
//芜湖乐善网络科技有限公司  乐善（惠游）
//武汉邦万科技发展有限公司  乐嗨嗨
//武汉多游科技有限公司
//武汉快玩科技有限公司  9665游戏
//武汉快游科技有限公司  悠游网
//武汉乐飞游科技有限公司
//武汉手盟网络科技有限公司
//武汉数游信息技术有限公司
//武汉扬程互联科技有限公司
//武汉游戏群科技有限公司
//武汉游侠精灵科技有限公司
//武汉有戏网络科技有限公司
//武汉掌游科技有限公司  安峰
//武汉爪游互娱科技有限公司
//西安玖毛网络科技有限公司
//西安闪游网络科技有限公司  骑士助手
//西华县龙城手游网络科技有限公司
//厦门巴掌科技有限公司
//厦门欢乐互娱网络科技有限公司  游动
//厦门市舜邦网络科技有限公司
//厦门说玩互娱科技有限公司
//厦门仙侠网络股份有限公司
//厦门小皮网络有限公司
//厦门鹰游网络科技有限公司
//厦门征游网络科技有限公司 43997
//新疆乱码网络科技有限公司
//炫彩互动网络科技有限公司
//烟台风暴网络科技有限公司
//鹰潭网源科技有限公司
//浙江欢游网络科技有限公司
//浙江游菜花网络科技有限公司 找乐
//浙江游菜花网络科技有限公司 CPS
//浙江游菜花网络科技有限公司 联运
//重庆玖度科技有限公司  爱上游戏
//重庆炎尚网络科技有限公司
//重庆宴门网络科技有限公司`
//
//	companies := []string{}
//	for _, v := range strings.Split(s, "\n") {
//		companies = append(companies, strings.Split(v, " ")[0])
//	}
//
//	// 得到公司id
//	companyids := []string{}
//
//	for _, v := range companies {
//		c := Company{}
//		orm.NewOrm().QueryTable("company").Filter("name", v).One(&c, "id")
//		companyids = append(companyids, strconv.Itoa(c.Id))
//	}
//
//	companyIdStr := strings.Join(companyids, "\n")
//
//	log.Info("test", companyIdStr)
//}

// 设置对账单与汇款单的偏移量，让数据统计正确
func TestFixChannelAmount(t *testing.T) {
	//s := `480|武汉爪游互娱科技有限公司|3764.90||1`
	//s := `公司id|公司名|未回款金额|预付款|主体`

	s := `238|广州游帮信息科技有限公司|556.14||1
209|广州千骐动漫有限公司|1,028.38||1
253|深圳市迅雷网络技术有限公司|1,183.14||1
471|上海草花互动网络科技有限公司|1,252.59||1
463|广东天宸网络科技有限公司|1,353.75||1
469|北京爱奇艺科技有限公司|1,589.35||1
113|厦门征游网络科技有限公司|2,015.90||1
244|上海瑛麒动漫科技有限公司|3,338.78||1
480|武汉爪游互娱科技有限公司|3,764.90||1
170|深圳泡椒思志信息技术有限公司|4,428.34||1
290|深圳豪邦网络有限公司|4,525.30||1
87|上海聚力传媒技术有限公司|4,825.05||1
165|北京聚乐网络科技有限公司|5,153.47||1
140|上饶市奇天乐地科技有限公司|5,549.55||1
56|苏州蜗牛数字科技股份有限公司|5,869.26||1
111|广州市久邦数码科技有限公司|7,225.13||1
123|北京广狐信息技术有限公司|7,834.41||1
239|炫彩互动网络科技有限公司|8,386.13||1
477|广州想玩信息科技有限公司|8,724.93||1
205|广州沙巴克网络科技有限公司|11,554.77||1
191|上海邑世网络科技有限公司|12,557.07||1
25|木蚂蚁（北京）科技有限公司|12,633.53||1
534|武汉手盟网络科技有限公司|14,365.00||1
161|广州游玩网络科技有限公司|15,506.37||1
249|南京直立行走网络科技有限公司|16,188.95||1
144|上海绿岸网络科技股份有限公司|16,844.55||1
12|北京力天无限网络技术有限公司|18,367.93||1
134|福建游龙网络科技有限公司|18,745.55||1
207|成都闲玩网络科技有限公司|19,070.44||1
256|贵州指趣网络科技有限公司|21,632.91||1
329|湖南草花互动网络科技有限公司|22,523.55||1
194|上海掌越科技有限公司|31,118.63||1
166|福州靠谱网络有限公司|33,107.74||1
122|海南汉风科技有限公司|40,277.33||1
91|上海晋昶网络科技有限公司|44,168.03||1
121|北京海誉动想科技股份有限公司|46,545.84||1
118|杭州宝亮网络科技有限公司|52,564.63||1
262|鹰潭网源科技有限公司|52,663.73||1
197|成都盟宝互动科技有限公司|56,432.60||1
240|浙江欢游网络科技有限公司|65,975.89||1
147|新疆乱码网络科技有限公司|73,902.39||1
214|努比亚技术有限公司|82,896.53||1
2|天津安果科技有限公司|89,069.99||1
177|上海弘贯网络科技有限公司|234899.57||1
237|北京勇士互娱科技有限公司|98,957.71||1
148|北京多点在线科技有限公司|107,877.45||1
274|南京规贝软件科技有限公司|120,644.11||1
479|上海娱嘉网络科技有限公司|124,261.50||1
424|绵阳友游网络科技有限公司|126,479.30||1
141|武汉快玩科技有限公司|168,277.68||1
257|上海都玩网络科技有限公司|177,470.93||1
95|深圳市拇指游玩科技有限公司|298,407.10||1
129|北京趣玩互动信息技术有限公司|319,917.10||1
137|重庆玖度科技有限公司|414,484.52||1
10|北京拇指玩科技有限公司|433,867.49||1
230|四三九九网络股份有限公司|465,277.68||1
119|天行（上海）网络科技有限公司|474,140.33||1
138|上海云群实业网络科技有限公司|478,097.65||1
213|福建乐游网络科技有限公司|680,772.85||1
195|广州道先网络科技有限公司|866,799.54||1
133|北京掌中悦动科技有限公司|1,042,374.89||1
153|上海游戏多网络科技股份有限公司|1,268,415.32||1
158|广州水煮信息科技有限公司|1,418,302.19||1
163|北京机锋科技有限公司|1,769,665.73||1
116|武汉掌游科技有限公司|2,369,749.07||1
196|广东秒创网络科技股份有限公司|3,584,181.66||1
127|北京蘑菇玩科技有限公司|3,221,208.50||1
128|广州小朋网络科技有限公司|1437925.64||1
484|宿迁楚风互娱科技有限公司|7,733,481.51||1
221|厦门市舜邦网络科技有限公司|220.40||2
219|上海巽力互联网科技有限公司|235.13||2
292|深圳市花生科技有限公司|399.95||2
180|北京掌汇天下科技有限公司|432.73||2
156|重庆炎尚网络科技有限公司|1,119.30||2
112|视娱（天津）网络科技有限公司|1,346.15||2
12|北京力天无限网络技术有限公司|1,355.78||2
148|北京多点在线科技有限公司|1,444.95||2
226|湖南奇葩乐游网络科技有限公司|1,630.20||2
272|杭州顺网宇酷科技有限公司|1,646.83||2
128|广州小朋网络科技有限公司|2,185.95||2
181|北京当乐信息技术有限公司|2,239.14||2
290|深圳豪邦网络有限公司|4,702.29||2
165|北京聚乐网络科技有限公司|5,272.47||2
81|厦门仙侠网络科技有限公司|7,070.81||2
25|木蚂蚁（北京）科技有限公司|12,554.26||2
119|天行（上海）网络科技有限公司|12,796.60||2
137|重庆玖度科技有限公司|13,690.64||2
500|长沙超神网络有限公司|17,262.14||2
487|北京手游端享科技有限公司|18,204.85||2
240|浙江欢游网络科技有限公司|20,515.25||2
141|武汉快玩科技有限公司|23,089.28||2
201|北京乐盟互动科技有限公司|24,899.50||2
245|重庆宴门网络科技有限公司|38,930.87||2
613|海马云（天津）信息技术有限公司|155,895.94||2
126|武汉邦万科技发展有限公司|57,856.89||2
159|武汉数游信息技术有限公司|57,443.27||2
499|北京鹏泰互动广告有限公司||1,247.50|2
190|广东安久科技有限公司||97,352.93|2
256|贵州指趣网络科技有限公司||95,262.73|2
504|鹤壁逍遥游网络科技有限公司||1,542.38|2
541|湖南天宇游网络科技有限公司||35,765.58|2
410|湖南炎夏网络科技有限公司||9,770.72|2
558|会理县云发现网络科技有限公司||8,677.12|2
424|绵阳友游网络科技有限公司||262,275.74|2
507|上海悦俊信息科技有限公司||10,000.00|2
506|郑州市中奥信息技术有限公司||29,233.86|2
142|厦门欢乐互娱网络科技有限公司||430467.35|1
217|厦门鹰游网络科技有限公司||78782.91|1
226|湖南奇葩乐游网络科技有限公司||27954.18|1
294|杭州宝莲灯科技有限公司||31869.01|1
258|广州火烈鸟网络科技有限公司|17047.55||1
286|广州裕际网络科技有限公司||18096.13|1
287|龙川奇点网络科技有限公司||15647.67|1
146|武汉扬程互联科技有限公司|15423.13||1`

	type ChannelAmount struct {
		companyId int
		amount    float64
		bodyMy    int
	}

	var cs []ChannelAmount
	for _, line := range strings.Split(s, "\n") {
		ds := strings.Split(line, "|")
		companyId, _ := strconv.Atoi(ds[0])
		var amount float64 = 0
		if ds[2] != "" {
			a := strings.Replace(ds[2], ",", "", -1)
			amount, _ = strconv.ParseFloat(a, 64)
			amount = - amount
		} else if ds[3] != "" {
			a := strings.Replace(ds[3], ",", "", -1)
			amount, _ = strconv.ParseFloat(a, 64)
		}
		bodyMy, _ := strconv.Atoi(ds[4])

		cs = append(cs, ChannelAmount{
			companyId: companyId,
			amount:    amount,
			bodyMy:    bodyMy,
		})
	}

	// 全部汇款单
	sql := "SELECT body_my,remit_company_id,SUM(amount) as amount FROM `remit_down_account` WHERE remit_time<1497499200  GROUP BY body_my,remit_company_id"
	var allRemit []orm.Params
	orm.NewOrm().Raw(sql).Values(&allRemit)
	allRemitAmount := map[string]ChannelAmount{}
	for _, v := range allRemit {
		bodyMy, _ := util.Interface2Int(v["body_my"], false)
		companyId, _ := util.Interface2Int(v["remit_company_id"], false)
		amount, _ := util.Interface2Float(v["amount"], false)
		allRemitAmount[fmt.Sprintf("%d %d", bodyMy, companyId)] = ChannelAmount{
			amount:    amount,
			companyId: int(companyId),
			bodyMy:    int(bodyMy),
		}
	}

	// 全部对账单（状态30/100）
	sql = "SELECT body_my,remit_company_id,SUM(amount_payable) as amount FROM `verify_channel` WHERE date <='2017-06' AND `status` >20  GROUP BY body_my,remit_company_id;"
	var allVerify []orm.Params
	orm.NewOrm().Raw(sql).Values(&allVerify)
	allVerifyAmount := map[string]ChannelAmount{}
	for _, v := range allVerify {
		bodyMy, _ := util.Interface2Int(v["body_my"], false)
		companyId, _ := util.Interface2Int(v["remit_company_id"], false)
		amount, _ := util.Interface2Float(v["amount"], false)
		allVerifyAmount[fmt.Sprintf("%d %d", bodyMy, companyId)] = ChannelAmount{
			amount:    amount,
			companyId: int(companyId),
			bodyMy:    int(bodyMy),
		}
	}

	// 所有在对账单和回款单出现的公司
	// 现在有这样一个逻辑：不在上表上面的公司 都是已经全部回款了的，要将偏移量置为0
	allCompanyAndBody := map[string]bool{}
	for k := range allVerifyAmount {
		allCompanyAndBody[k] = false
	}
	for k := range allRemitAmount {
		allCompanyAndBody[k] = false
	}

	hasAmountCompany := map[string]ChannelAmount{}
	for _, channelAmount := range cs {
		if channelAmount.companyId == 0 {
			continue
		}

		k := fmt.Sprintf("%d %d", channelAmount.bodyMy, channelAmount.companyId)
		allCompanyAndBody[k] = false
		hasAmountCompany[k] = channelAmount
	}

	for key := range allCompanyAndBody {
		bodyMy := 0
		companyId := 0

		channelAmount, has := hasAmountCompany[key]
		if has {
			bodyMy = channelAmount.bodyMy
			companyId = channelAmount.companyId
		}

		// 找到汇款单
		amountRemit, has := allRemitAmount[key]
		if has {
			bodyMy = amountRemit.bodyMy
			companyId = amountRemit.companyId
		}
		// 找到对账单
		amountVerify, has := allVerifyAmount[key]
		if has {
			bodyMy = amountVerify.bodyMy
			companyId = amountVerify.companyId
		}

		// 回款单 - 对账单 = 余额 + 偏移量
		// 偏移量 = 回款单 - 对账单 - 余额
		//println(amountRemit.amount)
		//println(amountVerify.amount)
		//println(channelAmount.amount)
		beego.Warning(amountRemit, amountVerify, channelAmount)

		amountOffset := amountRemit.amount - amountVerify.amount - channelAmount.amount
		beego.Warning(amountRemit.amount - amountVerify.amount - channelAmount.amount)

		// 把偏移量 和 真正的余额存起来
		err := UpdateOrCreateRemitPre(bodyMy, companyId, channelAmount.amount, amountOffset)
		if err != nil {
			log.Error("test", err)
		}
	}

	log.Info("test", cs)
}

// 同步pre_order表
func TestUpdatePreVerifyCpFromOrder(t *testing.T) {
	affCount, _ := UpdatePreVerifyCpFromOrder(nil, "2018-01")
	fmt.Printf("affCount:%v\n", affCount)
	//UpdatePreVerifyCpFromOrder(nil, "2017-12")
	//UpdatePreVerifyChannelFromOrder([]interface{}{"6816"}, "2017-12")

	//for i := -18; i < 0; i ++ {
	//	timeString := time.Now().AddDate(0, i, 0).Format("2006-01")
	//	fmt.Printf("titmeString:%v\n", timeString)
	//	affCount, _ := UpdatePreVerifyCpFromOrder([]interface{}{57}, timeString)
	//	fmt.Printf("affCount:%v\n", affCount)
	//}
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
		"kuaifazs", "10.8.230.17",
		"3308", "work_together_online")
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
	//	"123456", "10.8.230.17",
	//	"3308", "work_together")
	fmt.Printf("link:%v", link)
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
