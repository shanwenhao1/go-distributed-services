package main

import (
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Hello struct {
	Who string
}

type HelloActor struct{}

// Receive is the handle of request
func (state *HelloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Hello:
		fmt.Println(fmt.Sprintf("Hello %s", msg.Who))
	}
}

func main() {
	props := actor.PropsFromProducer(func() actor.Actor { return &HelloActor{} })
	rootContext := actor.EmptyRootContext
	pid := rootContext.Spawn(props)
	rootContext.Send(pid, &Hello{Who: "Roger"})
	console.ReadLine()
}
