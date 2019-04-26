/*
	通用返回, 请求处理模块
*/
package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-distributed-services/infra/enum"
	"go-distributed-services/infra/log"
	"io/ioutil"
	"net/http"
)

type Req interface {
	GetRequestData() interface{}
}

// 通用消息返回模型
type ParamModel struct {
	ErrorCode int
	ErrorMsg  interface{}
	Obj       interface{}
	OtherData interface{}
}

// 通用返回处理函数
func (par ParamModel) CommonResponse(ctx *gin.Context) {
	if par.ErrorCode == 0 {
		GetSuccessRJM(par.Obj).ResponseData(ctx)
	} else {
		GetDefaultRJM(par.ErrorCode).ResponseData(ctx)
	}
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

func (request RequestJsonModel) GetRequestData(ctx *gin.Context, rjm interface{}) interface{} {
	var reqData RequestJsonModel
	req := ctx.Request
	addr := req.Header.Get("X-Real-IP") // 获取真实发出请求的客户端IP
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For") // 获取IP(包含代理IP）
		if addr == "" {
			addr = req.RemoteAddr
		}
	}
	log.LogWithTag(log.INFO, log.ReqParse, "Request %s for %s", req.URL.Path, addr)
	dataS, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.LogWithTag(log.ERROR, log.ReqParse, "%w : %w", "Gin Read Body Error", err)
	}
	log.LogWithTag(log.INFO, log.ReqParse, "%v : %v", "The Request Body is", string(dataS))
	err = json.Unmarshal(dataS, rjm)
	if err != nil {
		log.LogWithTag(log.ERROR, log.ReqParse, "%v : %v", "Convert Body To Json Failed", err)
	}
	json.Unmarshal(dataS, &reqData)
	// you can do something with request obj
	if err != nil {
		GetDefaultRJM().ResponseData(ctx)
		return nil
	} else {
		return rjm
	}
}

//响应数据模型
type ResponseJsonModel struct {
	Obj       interface{} `json:"obj"`       // 内容
	ErrorCode int         `json:"errorCode"` // 编码
	Token     interface{} `json:"token"`     // token
	ErrorMsg  interface{} `json:"errorMsg"`  // 消息
}

/*
	响应函数
*/
func (response ResponseJsonModel) ResponseData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response)
}

//获取默认返回消息模型
func GetDefaultRJM(code ...int) ResponseJsonModel {
	if len(code) > 0 {
		return ResponseJsonModel{ErrorCode: code[0], ErrorMsg: enum.CodeMap[code[0]]}
	} else {
		return ResponseJsonModel{ErrorCode: enum.OPERATE_FAILED, ErrorMsg: enum.CodeMap[enum.OPERATE_FAILED]}
	}
}

//获取成功返回消息模型
func GetSuccessRJM(params ...interface{}) ResponseJsonModel {
	if len(params) == 1 {
		return ResponseJsonModel{ErrorCode: enum.OPERATE_SUCCESS, ErrorMsg: enum.CodeMap[enum.OPERATE_SUCCESS], Obj: params[0]}
	}
	if len(params) == 2 {
		return ResponseJsonModel{ErrorCode: enum.OPERATE_SUCCESS, ErrorMsg: enum.CodeMap[enum.OPERATE_SUCCESS], Obj: params[0], Token: params[1]}
	}
	return ResponseJsonModel{ErrorCode: enum.OPERATE_SUCCESS, ErrorMsg: enum.CodeMap[enum.OPERATE_SUCCESS]}
}
