package main

import (


	"github.com/zhangweilun/goweb/config"
	"github.com/zhangweilun/tinyweb"
)

/**
* 
* @author willian
* @created 2017-02-22 14:17
* @email 18702515157@163.com  
**/
func main() {
	config.Init()

	web := tinyweb.Classic()
	web.UseHandler(api.Api())
	web.Run(common.Api_port)
}
