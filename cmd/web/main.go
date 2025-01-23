package main

import (
	"mediamaid/app/global/variable"
	_ "mediamaid/bootstrap"
	"mediamaid/routers"
)

// 这里可以存放后端路由（例如后台管理系统）
func main() {
	router := routers.InitWebRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Web.Port"))
}
