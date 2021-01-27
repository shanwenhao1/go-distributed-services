package action

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gin-gonic/gin"
	"go-distributed-services/domain/model"
	"go-distributed-services/domain/services"
	"go-distributed-services/infra/log"
	"time"
)

// user action request parameter
type UserJsonModel struct {
	model.RequestJsonModel
	Obj model.User `json:"obj"`
}

// register request
func (this UserJsonModel) Register(c *gin.Context) {
	// 获取请求参数, 可以考虑在此添加中间件
	rjm := this.GetRequestData(c, &this.RequestJsonModel)
	if rjm != nil {
		jsonModel := *rjm.(*model.RequestJsonModel)
		// 使用actor消息驱动
		props := actor.PropsFromProducer(func() actor.Actor { return &service.UserActor{} })
		actContext := actor.RootContext{}
		pid := actContext.Spawn(props)
		handle := service.UserHandle{
			ReqJson:      jsonModel,
			HandleChoice: service.Register,
		}
		result, err := actContext.RequestFuture(pid, handle, time.Second*2).Result()
		if err != nil {
			log.LogWithTag(log.ERROR, log.InitSer, err.Error())
			model.GetDefaultRJM().ResponseData(c)
		}
		resp := result.(model.ParamModel)
		resp.CommonResponse(c)
	}
}

// login request
func (this UserJsonModel) Login(c *gin.Context) {
	// 获取请求参数, 可以考虑在此添加中间件
	rjm := this.GetRequestData(c, &this.RequestJsonModel)
	if rjm != nil {
		jsonModel := *rjm.(*model.RequestJsonModel)
		// 使用actor消息驱动
		props := actor.PropsFromProducer(func() actor.Actor { return &service.UserActor{} })
		actContext := actor.RootContext{}
		pid := actContext.Spawn(props)
		handle := service.UserHandle{
			ReqJson:      jsonModel,
			HandleChoice: service.Login,
		}
		result, err := actContext.RequestFuture(pid, handle, time.Second*2).Result()
		if err != nil {
			log.LogWithTag(log.ERROR, log.InitSer, err.Error())
			model.GetDefaultRJM().ResponseData(c)
		}
		resp := result.(model.ParamModel)
		resp.CommonResponse(c)
	}
}
