package test

import (
	"flag"
	"fmt"
	"go-distributed-services/infra/log"
	"testing"
)

func MockLog(t testing.T) {

}

func TestMain(m *testing.M) {
	fmt.Println("Test Begin")
	flag.Parse()
	log.InitializedLog4go("../config/logTest.xml")
	m.Run()
	fmt.Println("Test End")
}
