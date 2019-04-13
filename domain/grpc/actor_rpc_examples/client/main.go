package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"go-distributed-services/domain/grpc/actor_rpc_examples/messages"
)

func main() {
	timeout := 5 * time.Second
	remote.Start("127.0.0.1:8081")
	pidResp, err := remote.SpawnNamed("127.0.0.1:8080", "remote", "hello", timeout)
	if err != nil {
		panic(err)
	}
	pid := pidResp.Pid
	fmt.Println("----------", pid)
	res, err := actor.EmptyRootContext.RequestFuture(pid, &messages.HelloRequest{}, timeout).Result()
	if err != nil {
		panic(err)
	}
	response := res.(*messages.HelloResponse)
	fmt.Printf("Response from remote %v", response.Message)

	console.ReadLine()
}
