package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-distributed-services/infra/log"
	"strings"
)

type RConfig struct {
	MPrefix     string `xml:"prefix"`
	MPort       string `xml:"port"`
	MEnv        string `xml:"serverModel"`
	MUploadPath string `xml:"uploadPath"`
	MSsl        bool   `xml:"useSSL"`
}

var ConfigDataS RConfig

// 网络框架初始化(注册回调函数供加载路由)
func Init() {
	handlerMap := make(map[string]gin.HandlerFunc)
	// 设置为线上环境
	gin.SetMode(gin.ReleaseMode)
	//router := gin.Default()
	router := gin.New()
	Router(handlerMap)
	log.LogWithTag(log.INFO, log.InitSer, "http网络框架初始化完成[%s][%s]", ConfigDataS.MPort, ConfigDataS.MEnv)
	for patten, handle := range handlerMap {
		// 对特殊的路由处理, 其余一律采用POST方法
		if strings.Contains(patten, "Upload") {
			log.LogWithTag(log.INFO, log.InitSer, "文件下载插件注册完成[%s]", patten)
			//router.Any("/test/Upload/", handle)
		} else {
			router.POST(ConfigDataS.MPrefix+patten, handle)
		}
	}
	// 另一种注册路由方式
	router.Static("/test/Upload/", "./Upload")
	// router.RunTLS使用HTTPS加密连接(需生成ssl key), router.Run使用http连接
	if ConfigDataS.MSsl {
		router.RunTLS(ConfigDataS.MPort, "keys/my.pem", "keys/my.key")
	} else {
		fmt.Println("xxxxxxxxxxx Server Run On", ConfigDataS.MPort, "...")
		router.Run(ConfigDataS.MPort)
	}
}
