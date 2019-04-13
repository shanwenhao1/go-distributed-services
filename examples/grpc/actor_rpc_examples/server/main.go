package main

import (
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"go-distributed-services/examples/grpc/actor_rpc_examples/messages"
)

type MyActor struct{}

func (*MyActor) Receive(context actor.Context) {
	fmt.Println("----------receive request")
	switch msg := context.Message().(type) {
	case *messages.Echo:
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	default:
		fmt.Println("=============", msg)
	}
}

func main() {
	remote.Start("localhost:8091")

	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", actor.PropsFromProducer(func() actor.Actor { return &MyActor{} }))
	console.ReadLine()
}
