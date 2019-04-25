package example

import (
	"github.com/golang/mock/gomock"
	"go-distributed-services/test/example/spider"
	"testing"
)

// 使用
// mockgen -destination test/example/spider/mock_spider.go -package spider go-distributed-services/test/example/spider Spider
// generate mock_spider模拟spider接口的实现
func TestGetGoVersion(t *testing.T) {
	// 使用mock生成预期数据进行测试
	mockCtl := gomock.NewController(t)
	mockSpi := spider.NewMockSpider(mockCtl)
	mockSpi.EXPECT().GetBody().Return("go1.8.9")
	v := spider.GetGoVersion(mockSpi)
	if v != "go1.8.9" {
		t.Error("Get wrong version %s", v)
	}
}

// 详情: https://www.jianshu.com/p/598a11bbdafb
