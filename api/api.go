package api

import (
	"github.com/julienschmidt/httprouter"

	"github.com/zhangweilun/goweb/api/font"
)

/**
*
* @author willian
* @created 2017-01-27 14:48
* @email 18702515157@163.com
**/

func Api() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", font.Index)
	router.GET("/proxy", font.GetProxy)
	router.DELETE("/delete/:id",font.Delete)
	router.GET("/count",font.Count)
	return router
}
