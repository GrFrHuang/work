package test

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"strings"
	"testing"

	"net/http"
	"io/ioutil"
	"github.com/astaxie/beego"
)

func init() {
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kftest",
	//	"123456", "10.8.230.17",
	//	"3308", "work_together")
	////link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
	////	"kuaifazs", "10.8.230.17",
	////	"3308", "work_together_online")
	//orm.RegisterDataBase("default", "mysql", link)
	////orm.Debug = true
	//
	//_, file, _, _ := runtime.Caller(1)
	//apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	//beego.TestBeegoInit(apppath)
}

func TestTime(t *testing.T) {
	/*token := "94dcf500e220f02ed775abe4ea9c87a0"
	cps := beego.AppConfig.String("cps")
	uri := cps + "/v1/message_scrollbar/setmessaget?token=" + token
	_, err := url.Parse(uri)*/
	//req, err := http.Get("http://www.baidu.com")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer req.Body.Close()
	//data,err:=ioutil.ReadAll(req.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(data))


	token := "943f0348775d794fc480dd38c04d14ff"
	cps := beego.AppConfig.String("cps")
	uri := cps+"/v1/message_scrollbar/setmessaget?token="+token
	//_, err = url.Parse(uri)

	req, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}




////func TestTime(t *testing.T) {
//	tim, err := time.Parse("2006-01-02", "2030-01-02")
//	if err != nil {
//		return
//	}
//	x := int32(tim.Unix())
//
//	log.Print(x)
//}

// fix 渠道商(channel_company)的回款主体 (还能添加公司)
// company表找公司
//  找到：得到公司id
//  没找到：新建公司
// 找出 channel_code
//    不存在：
//			生成一条channel_company (company_id,remit_company_id)
//    存在：
//      重置remit_company_id
func TestFixChannelRemitCompany(t *testing.T) {

	s := `360	360手机助手	北京世界星辉科技有限责任公司
kuaiyong	快用	天津安果科技有限公司
dl	当乐	北京当乐信息技术有限公司
vivo	VIVO	广东天宸网络科技有限公司
4399	4399	四三九九网络股份有限公司
wandoujia	豌豆荚	北京卓易讯畅科技有限公司
huawei	华为	华为软件技术有限公司
muzhiwan	拇指玩	北京拇指玩科技有限公司
jolo	聚乐	北京聚乐网络科技有限公司
sogou	搜狗	北京搜狗网络技术有限公司
oppo	oppo	东莞市讯怡电子科技有限公司
youlong	游龙工会	福建游龙网络科技有限公司
itools	itools	深圳市创想天空科技股份有限公司
49you	世加游戏	上海晋昶网络科技有限公司
xunlei	迅雷	深圳市迅雷网络技术有限公司
umi	有米偶玩	淮安有米信息科技有限公司
hucn	乐非凡	广州云翎网络科技有限公司
pickle	泡椒	深圳泡椒思志信息技术有限公司
seahorse	海马	北京海誉动想科技股份有限公司/大胜数娱（天津）科技有限公司
mumayi	木蚂蚁	木蚂蚁（北京）科技有限公司
2324	3G门户	广州市久邦数码科技有限公司
gfan	机锋	北京机锋科技有限公司
i4	爱思	深圳市为爱普信息技术有限公司
greencoast	绿岸	上海绿岸网络科技股份有限公司
pps	PPS	北京爱奇艺科技有限公司
uu	悠悠村	上海优扬新媒信息技术有限公司
qixiazi	七匣子	杭州快定网络股份有限公司
youxiqun	游戏群	武汉游戏群科技有限公司
boyakenuo	三星博雅科诺	北京博雅科诺信息技术有限公司
much	摩奇手机(鱼丸互动)	深圳市鱼丸互动科技有限公司
caohua	51手游社区_草花	上海草花互动网络科技有限公司
9133	手游帮	广州游帮信息科技有限公司
appchina	应用汇	上海范特西网络科技有限公司
anfan	安锋	武汉掌游科技有限公司
8868	8868	广州水煮信息科技有限公司
qipa	树海(试玩_奇葩)	湖南奇葩乐游网络科技有限公司
zs	值尚	深圳市值尚互动科技有限公司
haima	海马安卓	北京海誉动想科技股份有限公司/大胜数娱（天津）科技有限公司
07073sy	数游	武汉数游信息技术有限公司
snail	免商店	苏州蜗牛数字科技股份有限公司
lewanduo	乐玩多平台	上海爱去玩科技有限公司
91wan	91WAN	广州维动网络科技有限公司
44755	奇天乐地	上饶市奇天乐地科技有限公司
guopan	果盘	广州火烈鸟网络科技有限公司
xxzhushou	叉叉助手(果盘)	广州火烈鸟网络科技有限公司
yayawan	丫丫玩	广州千骐动漫有限公司
6816	职内科技	厦门仙侠网络股份有限公司
93pk	都玩93	上海都玩网络科技有限公司
giant	手游咖啡	上海鸿游信息技术有限公司
weibo	新浪游戏	微梦创科网络科技（中国）有限公司
lepay	乐视	乐视移动智能信息技术（北京）有限公司
7k7ksy	7k7k手游	北京奇客创想信息技术有限公司
hupu	虎扑	深圳市星耀互动科技有限公司
qxz	快定网络	杭州快定网络股份有限公司
23youxi	爱上游戏（玖度）	重庆玖度科技有限公司
kaopu	靠谱	福州靠谱网络科技有限公司
snssdk	今日头条	北京字节跳动科技有限公司
egameddai	爱游戏短代	炫彩互动网络科技有限公司
shunwang	顺网	杭州顺网宇酷科技有限公司
biwan	必玩	上海都玩网络科技有限公司
lehihi	乐嗨嗨	武汉邦万科技发展有限公司
xiongbing	雄兵	广州雄兵网络科技有限公司
leba	乐8	鹰潭网源科技有限公司
dyoo	点优	广州点优网络科技有限公司
mgw	芒果玩	广州唐腾信息技术有限公司
benshouji	笨手机	天行（上海）网络科技有限公司
ccplay	虫虫助手	上海掌越科技有限公司
lb	江苏猎宝	南京规贝软件科技有限公司
huli	狐狸助手	北京广狐信息技术有限公司
77l	齐齐乐（发发）	江西巨网科技股份有限公司
linnyou	麟游SDK	成都盟宝互动科技有限公司
17sy	17手游吧	广州天势网络科技有限公司
yoyou	优游网(快游)	武汉快游科技有限公司
nox	夜神	北京多点在线科技有限公司
igamecool	找乐助手	浙江游菜花网络科技有限公司
tt	谊游（TT语音）	广州沙巴克网络科技有限公司
cuudoo	迷思(酷嘟)	广州想玩信息科技有限公司
7723	7723	厦门巴掌互动科技有限公司
anqu	安趣	北京安趣科技股份有限公司
muzhiplat	拇指游玩	深圳市拇指游玩科技有限公司
ququ	趣趣	北京趣趣网络科技有限公司
sguo	松果	福建创意嘉和软件有限公司
pyw	朋友玩	广州小朋网络科技有限公司
wqwan	我去玩	深圳飓风伟业网络有限公司
erhu	二狐游戏	厦门市舜邦网络科技有限公司
qianbao	钱宝	上海冰穹网络科技有限公司
dadayou	大大游(龙翔)	江苏龙翔网络科技有限公司
nifenglin	逆风鳞	广州逆风鳞科技有限公司
papa	啪啪游戏厅	浙江欢游网络科技有限公司
521play	521手游	上海娱嘉网络科技有限公司
xianwan	闲玩	成都闲玩网络科技有限公司
7u	7U（梦娱）	河北梦娱星空网络科技有限公司
qmzs	全民助手	北京趣玩互动信息技术有限公司
1881wan	1881玩	泉州市海悦网络科技有限公司
youxia	游侠	武汉游侠精灵科技有限公司
zhuayou	抓游（安峰小号）	武汉爪游互娱科技有限公司
xgame	YY手游宝	广州华多网络科技有限公司
damai	大麦助手	广州道先网络科技有限公司
youxiduo	游戏多	上海游戏多网络科技股份有限公司
tanwan	贪玩	江西贪玩信息技术有限公司
shangbar	乐盟	北京乐盟互动科技有限公司
3456wan	3456玩（邑世游戏）	上海邑世网络科技有限公司
51wan	51玩	北京新娱兄弟网络科技有限公司
shuowan	说玩	厦门说玩互娱科技有限公司
kuangwan	狂玩	广东秒创网络科技股份有限公司/杭州掌游科技有限公司
9665	9665游戏	武汉快玩科技有限公司
baolechufeng	宝乐楚风	宿迁楚风互娱科技有限公司/宿迁市宝乐网络科技有限公司
17168	豪邦网络	深圳豪邦网络有限公司
weme	微米	深圳微米动力科技有限公司
aiduoyou	武汉多游	武汉多游科技有限公司
shouyoudao	手游岛	深圳手游岛网络科技有限公司
moge	摩格	福建乐游网络科技有限公司
75757	惠游(75757)	芜湖乐善网络科技有限公司
100qu	百趣	南京时玳运成网络科技发展有限公司
nubia	努比亚	努比亚技术有限公司
6kwan	上海弘贯(6k玩)	上海弘贯网络科技有限公司
96788	游动96788	厦门欢乐互娱网络科技有限公司
fengbao	风暴	烟台风暴网络科技有限公司（风暴）
43997	征游	厦门征游网络科技有限公司
changba	唱吧	北京酷智科技有限公司
mgwyx	新芒果玩	广州唐腾信息技术有限公司
lcsygame	龙城手游	西华县龙城手游网络科技有限公司
2lyx	爱乐游戏	视娱（天津）网络科技有限公司
vqs	骑士助手	西安闪游网络科技有限公司
3011	玖毛	西安玖毛网络科技有限公司
lmyouxi	乱码游戏（科创网路）	新疆乱码网络科技有限公司
yiliu	一六游戏	广西一六游网络科技有限公司
ylsuper	游龙聚合SDK	北京掌中悦动科技有限公司
xintiao	心跳助手	北京手游端享科技有限公司
youkala	花生	深圳市花生科技有限公司
yshy	勇士互娱（酷游戏）	北京勇士互娱科技有限公司
papayou	啪啪游(厦游)	厦门厦游网络科技有限公司
youxifan	游戏Fan	广东安久科技有限公司
9quyx	九趣游戏	四川九趣网络科技有限公司
mushroom	怡扬(蘑菇玩)	北京蘑菇互娱科技有限公司/北京怡扬科技发展有限公司
iappsgame	爱应用IOS	南京直立行走网络科技有限公司
gm	怪猫	上海咕么信息科技有限公司
44937	44937手游天下	北京手游天下数字娱乐科技股份有限公司
yanmen	宴门手游	重庆宴门网络科技有限公司
x7sy	小七手游(尚米)	深圳尚米网络科技有限公司
qmdy	全民点游	深圳市全民点游科技有限公司
aiwan	爱玩	威海龙必达信息技术有限公司
xdsy	仙豆手游(幻动)	广州市幻动网络科技有限责任公司
shouyougou	手游狗(耀玩)	广州耀玩网络科技有限公司
putao	葡萄游戏	北京爱视游科技文化有限公司
lefeiyou	乐飞游(乐七)	武汉乐飞游科技有限公司/武汉扬程互联科技有限公司
zxgame	掌炫科技	广东秒创网络科技股份有限公司/杭州掌游科技有限公司
xiaoyao	逍遥	上海迈微软件科技有限公司
slsk	西柚(沙漏时空)	深圳西柚网络科技有限公司
manhuaren	漫画人	上海瑛麒动漫科技有限公司
baoliang	宝亮	杭州宝亮网络科技有限公司
yanshang	炎尚	重庆炎尚网络科技有限公司
baiyou	百游	上海巽力互联网科技有限公司
yaodian	遥望遥点	浙江游菜花网络科技有限公司
shunwan	顺玩	武汉生活范网络科技有限公司
hanfeng	汉风	海南汉风科技有限公司
dongdong	东东游戏(聚丰)	长沙七丽网络科技有限公司
ifeng	凤凰网	北京欢游天下科技有限公司
qiwan	奇顽游戏	北京奇顽科技有限公司
xd178	奥创(兄弟一起吧)	遂宁市奥创科技有限公司
dianjin	点金	福建点金网络科技有限公司
24235	嗨皮游戏(群云)	上海云群实业网络科技有限公司
youbao	游宝手游	成都游宝科技有限公司
yuwan	鱼丸互动	深圳市鱼丸互动科技有限公司
kingcheer	鲸旗天下	深圳鲸旗天下网络科技有限公司
lewanwx	乐玩无线	深圳乐玩无限科技有限公司
yeshen	新夜神	北京多点在线科技有限公司
xiaopi	小皮	厦门小皮网络有限公司
keshuo	柯烁(KK手游)	山东柯烁网络科技有限公司
taoshouyou	淘手游	贵州指趣网络科技有限公司
eagle	雄鹰游戏	厦门鹰游网络科技有限公司
aipu	爱谱	深圳市爱谱互娱可以有限公司
dmys	炎尚大麦	重庆炎尚网络科技有限公司
tuulang	螳螂游戏	北京中梦网络科技有限公司
opentide	三星鹏泰	北京鹏泰互动广告有限公司
ronbu	皮卡游戏	长沙超神网络有限公司
feichen	飞辰	杭州飞辰网络科技有限公司
mscgame	天趣游戏	福建天趣网络科技有限公司
cfwl	楚风网络	宿迁楚风互娱科技有限公司
regret	友游网络(后悔药)	绵阳友游网络科技有限公司
2yl	秒乐	广州妙乐网络科技有限公司
72g	72G(有戏)	武汉有戏网络科技有限公司
xyy	逍遥游	鹤壁逍遥游网络科技有限公司
qidianyule	奇点娱乐	龙川奇点网络科技有限公司
bell	贝尔游戏(广州闪趣)	广州闪趣网络科技有限公司
52tktk	TK游戏	广州裕际网络科技有限公司
shouyouzhu	手游猪(中奥)	中奥信息
lotuslantern	17玩吧(宝莲灯)	杭州宝莲灯科技有限公司
caoxie	草鞋	上海悦俊信息科技有限公司
11game
7676
baiguan
buer
customforkf
fox
grape
ihiyo
kudong
kunda
mkzoo
ninegame
peach
qix
qmgame
xinjigame
youliang
zhangwan`
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	ss := strings.Split(s, " ")

	ssT := []string{}
	for _, v := range ss {
		if v != "" {
			ssT = append(ssT, v)
		}
	}

	type X struct {
		code         string
		company_name string
	}

	xs := []X{}
	l := len(ssT)
	log.Println("test", l, ssT)
	for i := 0; i < l; i += 3 {
		xs = append(xs, X{
			code:         ssT[i],
			company_name: ssT[i+2],
		})

	}

	o := orm.NewOrm()
	for _, v := range xs {
		var companyId int64 = 0
		// 找公司
		sql1 := "SELECT * FROM `company` where `name` = ?"
		rs := []orm.Params{}
		name := strings.Split(v.company_name, "/")[0]
		o.Raw(sql1, name).Values(&rs)
		if len(rs) == 0 {
			log.Println("not found company: ", v.company_name)
			r, _ := o.Raw("INSERT into `company`(name,type) VALUES(?,3)", name).Exec()
			companyId, _ = r.LastInsertId()
		} else {
			companyId, _ = util.Interface2Int(rs[0]["id"], false)
		}

		if companyId == 0 {
			log.Println("found or created company error")
			continue
		}

		var companyId2 int64 = 0

		if strings.Contains(v.company_name, "/") {
			// 找公司
			sql12 := "SELECT * FROM `company` where `name` = ?"
			rs2 := []orm.Params{}
			name2 := strings.Split(v.company_name, "/")[1]
			o.Raw(sql12, name2).Values(&rs2)
			if len(rs2) == 0 {
				log.Println("not found company: ", v.company_name)
				r, _ := o.Raw("INSERT into `company`(name,type) VALUES(?,3)", name2).Exec()
				companyId2, _ = r.LastInsertId()
			} else {
				companyId2, _ = util.Interface2Int(rs2[0]["id"], false)
			}

		}

		// 找渠道商
		sql := "SELECT * FROM `channel_company` where channel_code = ?"
		rs = []orm.Params{}
		o.Raw(sql, v.code).Values(&rs)
		remitCompany := ""
		if companyId2 == 0 {
			remitCompany = fmt.Sprintf("[%d]", companyId)
		} else {
			remitCompany = fmt.Sprintf("[%d,%d]", companyId, companyId2)
		}
		if len(rs) == 0 {
			o.Raw("INSERT into `channel_company`(channel_code,company_id,remit_company) VALUES(?,?,?)",
				v.code,
				companyId,
				remitCompany,
			).Exec()

			log.Println("insert: ", v.code, companyId)
		} else {
			id, _ := util.Interface2Int(rs[0]["id"], false)
			o.Raw("UPDATE `channel_company` SET company_id= ? ,remit_company = ? WHERE id = ?",
				companyId,
				remitCompany,
				id,
			).Exec()

			log.Println("updated: ", v.code, companyId)

		}

	}

	log.Println(ssT)
}

func TestXX(t *testing.T) {
	r := 0
	// 在数组中,找出乘积最大的相邻元素组成的区间
	work([]int{3, 6, -2, 5, 7, 3}, 0, 5, &r)

	log.Print(r)

}

func product(inputArray []int, s, e int) int {
	if s >= e {
		return inputArray[s]
	} else {
		r := 1
		for i := s; i <= e; i++ {
			r *= inputArray[i]
		}
		return r
	}
}

func work(inputArray []int, s, e int, max *int) {
	if s >= e {
		return
	}

	m1 := product(inputArray, s, e)
	if m1 > *max {
		*max = m1
		log.Print(s, e)
	}

	work(inputArray, s+1, e, max)
	work(inputArray, s, e-1, max)

	return
}
