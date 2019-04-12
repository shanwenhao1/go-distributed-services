package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"go-distributed-services/domain/hello_example"
)

func main() {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &hello_example.HelloActor{}
	})
	rootContext := actor.EmptyRootContext
	pid := rootContext.Spawn(props)
	rootContext.Send(pid, &hello_example.Hello{Who: "Roger"})
	console.ReadLine()
}
