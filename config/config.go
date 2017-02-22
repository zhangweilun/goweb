package config

import (
	"os"
	"github.com/zhangweilun/logs"
	"github.com/astaxie/beego/orm"
)

/**
* 
* @author willian
* @created 2017-02-22 14:30
* @email 18702515157@163.com  
**/

func Init()  {
	work_space, _ := os.Getwd()
	logs.Init(work_space + "/logs")
	logs.Info("start ip proxy pool =============>>>>>>>>>>>>>>>>>>>")
	logs.Info("the logs directory is ",work_space+"/logs")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//root:root@/orm_test?
	orm.RegisterDataBase("default", "mysql", "admin:admin@tcp(master:3306)/crawler?charset=utf8&loc=Asia%2FShanghai", 30, 30)
	orm.RegisterModel(new(model.Proxy))
	orm.Debug = true
}