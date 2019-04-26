package service

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"go-distributed-services/domain/model"
	"go-distributed-services/infra/enum"
)

type LoginActor struct{}

func (login *LoginActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case model.RequestJsonModel:
		fmt.Println("-------------message", msg)
		resp := login.LoginH(msg)
		context.Send(context.Sender(), resp)
	}
}

// 登录操作
func (login LoginActor) LoginH(req model.RequestJsonModel) model.ParamModel {
	// do something(不关乎领域逻辑和业务的)
	fmt.Println("-------")
	// 登录逻辑
	return model.ParamModel{ErrorCode: enum.OPERATE_SUCCESS}
}
