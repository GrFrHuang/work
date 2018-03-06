package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// ALTER DATABASE db_name DEFAULT CHARACTER SET character_name
// ALTER TABLE table_name CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;
func Syncdb() {
	fmt.Println("database init is begin ...")
	var dataSource string
	db_host := beego.AppConfig.String("dev::mysqlurls")
	db_port := beego.AppConfig.String("dev::mysqlport")
	db_user := beego.AppConfig.String("dev::mysqluser")
	db_pass := beego.AppConfig.String("dev::mysqlpass")
	db_name := beego.AppConfig.String("dev::mysqldb")
	orm.RegisterModel(new(Game), new(Order), new(Channel))
	dataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	fmt.Println("url: ", dataSource)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dataSource)
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		fmt.Println(err)
	}
}
