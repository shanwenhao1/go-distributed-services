/*
	通用返回, 请求处理模块
*/
package model

// 通用消息返回模型
type ParamModel struct {
	ErrorCode int
	ErrorMsg  interface{}
	Obj       interface{}
	OtherData interface{}
}

//请求数据模型
type RequestJsonModel struct {
	AppId         string      `json:"appId"`
	Token         string      `json:"token"`
	Obj           interface{} `json:"obj"`
	ClientType    string      `json:"clientType"`
	Sign          string      `json:"sign"`
	TimeStamp     string      `json:"time_stamp"`
	ClientVersion string      `json:"clientVersion"`
}

//响应数据模型
type ResponseJsonModel struct {
	Obj       interface{} `json:"obj"`       // 内容
	ErrorCode int         `json:"errorCode"` // 编码
	Token     interface{} `json:"token"`     // token
	ErrorMsg  interface{} `json:"errorMsg"`  // 消息
}
