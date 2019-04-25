package interfaces

import (
	"encoding/xml"
	"go-distributed-services/infra/log"
	"go-distributed-services/interfaces/route"
	"io/ioutil"
	"runtime"
)

type Server struct{}

// 运行环境初始化(设置CPU核心数、读取server端口配置等)
func (this Server) InitializedSystem(path string) {
	dataS, rErr := ioutil.ReadFile(path)
	if rErr != nil {
		log.LogWithTag(log.ERROR, log.InitSer, "读取服务配置文件异常:[%v]", rErr)
		panic(rErr.Error())
	}
	configData := route.RConfig{}
	xErr := xml.Unmarshal(dataS, &configData)
	if xErr != nil {
		log.LogWithTag(log.ERROR, log.InitSer, "解析服务库配置文件异常:[%v]", xErr)
		panic(xErr.Error())
	}
	route.ConfigDataS = configData
	cc := runtime.NumCPU()
	// running in multi core
	runtime.GOMAXPROCS(cc)
	log.LogWithTag(log.INFO, log.InitSer, "运行环境初始化完成[处理器核心数:%d]", cc)
}

// 服务启动
func (this Server) Run() {
	route.Init()
}
