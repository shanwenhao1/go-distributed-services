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

// login request
func (this UserJsonModel) Login(c *gin.Context) {
	// 获取请求参数, 可以考虑在此添加中间件
	rjm := this.GetRequestData(c, &this.RequestJsonModel)
	if rjm != nil {
		jsonModel := *rjm.(*model.RequestJsonModel)
		// 使用actor消息驱动
		props := actor.PropsFromProducer(func() actor.Actor { return &service.LoginActor{} })
		actContext := actor.EmptyRootContext
		pid := actContext.Spawn(props)
		result, err := actContext.RequestFuture(pid, jsonModel, time.Second*2).Result()
		if err != nil {
			log.LogWithTag(log.ERROR, log.InitSer, err.Error())
			model.GetDefaultRJM().ResponseData(c)
		}
		resp := result.(model.ParamModel)
		resp.CommonResponse(c)
	}
}
