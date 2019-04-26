package action_test

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"go-distributed-services/domain/model"
	"go-distributed-services/domain/services"
	"go-distributed-services/infra/enum"
	"testing"
	"time"
)

func TestUserAction_Login(t *testing.T) {
	reqData := model.RequestJsonModel{}
	reqData.ClientType = "android"
	props := actor.PropsFromProducer(func() actor.Actor { return &service.LoginActor{} })
	actContext := actor.EmptyRootContext
	pid := actContext.Spawn(props)
	result, err := actContext.RequestFuture(pid, reqData, time.Second*2).Result()
	if err != nil {
		t.Error(err)
	}
	if result.(model.ParamModel).ErrorCode != enum.OPERATE_SUCCESS {
		t.Error("操作失败")
	}
}
