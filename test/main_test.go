package test

import (
	"flag"
	"fmt"
	"go-distributed-services/infra/log"
	"testing"
)

func MockLog(t testing.T) {

}

// use command `go test -v ./...` to run all test in package test
func TestMain(m *testing.M) {
	fmt.Println("Test Begin")
	flag.Parse()
	log.InitializedLog4go("../config/logTest.xml")
	m.Run()
	fmt.Println("Test End")
}
