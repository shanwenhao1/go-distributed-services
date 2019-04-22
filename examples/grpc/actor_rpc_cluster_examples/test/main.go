package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"go-distributed-services/examples/grpc/actor_rpc_cluster_examples/shared"
)

func main() {
	remote.Register("Hello", actor.PropsFromProducer(func() actor.Actor {
		return &shared.HelloActor{}
	}))
	console.ReadLine()
}
