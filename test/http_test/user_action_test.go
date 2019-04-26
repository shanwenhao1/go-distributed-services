package http_test

import (
	"go-distributed-services/infra/enum"
	"go-distributed-services/test"
	"testing"
)

func TestHttp_Login(t *testing.T) {
	data := map[string]interface{}{
		"login_type": test.LoginType,
		"user_name":  test.UserName,
		"password":   test.Password,
	}
	url := test.UrlIp + "/login"

	reqData := map[string]interface{}{
		"obj": data,
	}
	result, err := test.Request(url, reqData)
	if err != nil {
		t.Error(err)
	}
	if result.ErrorCode != enum.OPERATE_SUCCESS {
		t.Error("操作失败")
	}
}
