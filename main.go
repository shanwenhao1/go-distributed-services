package main

import (
	"go-distributed-services/infra/db"
	"go-distributed-services/infra/log"
	"go-distributed-services/interfaces"
)

func main() {
	ser := interfaces.Server{}
	// 日志初始化
	log.InitializedLog4go("config/log4go.xml")
	// 初始化数据库
	db.InitializedDataSource("config/dbConfig.xml")
	// 初始化启动
	ser.InitializedSystem("config/server_gin.xml")
	// 通过回调函数注册路由, 并启动服务
	ser.Run()
	//router := gin.Default()
	//s := &http.Server{
	//	Addr:           ":8080",
	//	Handler:        router,
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//s.ListenAndServe()
}
