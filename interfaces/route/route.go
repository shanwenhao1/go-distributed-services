package route

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gin-gonic/gin"
	"go-distributed-services/app/action"
)

type Route interface {
	RouteMessage(message interface{}, sender *actor.PID)
	SetRoutes(routes *actor.PIDSet)
	GetRoutes() *actor.PIDSet
}

// route register, register the handle function of  web request
func Router(handleMap map[string]gin.HandlerFunc) {
	userAction := new(action.UserJsonModel)

	// 添加login路由
	handleMap["register"] = userAction.Register
	handleMap["login"] = userAction.Login
}
