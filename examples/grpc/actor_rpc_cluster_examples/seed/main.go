package main

import (
	"log"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/hashicorp/consul/api"
	"go-distributed-services/examples/grpc/actor_rpc_cluster_examples/shared"
)

func main() {

	// this node knows about Hello kind
	remote.Register("Hello", actor.PropsFromProducer(func() actor.Actor {
		return &shared.HelloActor{}
	}))

	config := &api.Config{}
	config.Address = "192.168.1.89"

	cp, err := consul.NewWithConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("mycluster", "127.0.0.1:8080", cp)

	hello := shared.GetHelloGrain("MyGrain")

	res, err := hello.SayHello(&shared.HelloRequest{Name: "Roger"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message from grain: %v", res.Message)
	console.ReadLine()

	cluster.Shutdown(true)
}
