package main

import (
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"go-distributed-services/examples/grpc/actor_rpc_examples/messages"
	"time"
)

type MyActor struct {
	count int
}

func (state *MyActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *messages.Response:
		state.count++
		fmt.Println(state.count)
	}
}

func main() {
	remote.Start("localhost:8090")

	rootCtx := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor { return &MyActor{} })
	pid := rootCtx.Spawn(props)
	message := &messages.Echo{Message: "hej", Sender: pid}

	// this is to spawn remote actor we want to communicate with
	spawnResponse, err := remote.SpawnNamed("localhost:8091", "myactor", "hello", time.Second*10)
	if err != nil {
		panic(err)
	}

	// get spawned PID
	spawnedPID := spawnResponse.Pid
	for i := 0; i < 10; i++ {
		rootCtx.Send(spawnedPID, message)
	}

	console.ReadLine()
}
