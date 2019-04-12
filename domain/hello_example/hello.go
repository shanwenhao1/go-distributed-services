package hello_example

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Hello struct {
	Who string
}

type HelloActor struct{}

func (state *HelloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Hello:
		fmt.Println(fmt.Sprintf("Hello %s", msg.Who))
	}
}
