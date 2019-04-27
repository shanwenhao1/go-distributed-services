package service

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"go-distributed-services/domain/model"
	"go-distributed-services/infra/enum"
)

type UserHandle struct {
	ReqJson      model.RequestJsonModel
	HandleChoice int
}

const (
	Register = iota
	Login
)

type UserActor struct{}

// Actor执行模块, 注意事项:
// 1.命名必须为Receive(实际上是重载方法)
// 2.switch不支持default, 否则会报错
func (userAct *UserActor) Receive(context actor.Context) {
	switch handle := context.Message().(type) {
	case UserHandle:
		fmt.Println("==============", handle)
		num := handle.HandleChoice
		switch num {
		case Register:
			fmt.Println("-------------message", handle.ReqJson)
			resp := userAct.Register(handle.ReqJson)
			context.Send(context.Sender(), resp)
		case Login:
			resp := userAct.Login(handle.ReqJson)
			context.Send(context.Sender(), resp)
		}
	}
}

// 注册操作
func (userAct UserActor) Register(req model.RequestJsonModel) model.ParamModel {
	// do something(不关乎领域逻辑和业务的)
	fmt.Println("-------")
	// 注册逻辑
	return model.ParamModel{ErrorCode: enum.OPERATE_SUCCESS}
}

// 登录操作
func (userAct UserActor) Login(req model.RequestJsonModel) model.ParamModel {
	// do something(不关乎领域逻辑和业务的)
	fmt.Println("-------")
	// 登录逻辑
	return model.ParamModel{ErrorCode: enum.OPERATE_SUCCESS}
}
