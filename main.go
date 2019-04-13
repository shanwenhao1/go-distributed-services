package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"go-distributed-services/domain/actor_example"
)

func main() {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &actor_example.HelloActor{}
	})
	rootContext := actor.EmptyRootContext
	pid := rootContext.Spawn(props)
	rootContext.Send(pid, &actor_example.Hello{Who: "Roger"})
	console.ReadLine()
}
