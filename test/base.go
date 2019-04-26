package test

import (
	"encoding/json"
	"errors"
	"go-distributed-services/domain/model"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	UrlIp     = "http://127.0.0.1:8080/Test/:system"
	LoginType = "test"
	UserName  = "test_name"
	Password  = "password"
)

func Request(url string, reqData map[string]interface{}) (model.ParamModel, error) {
	var result model.ParamModel
	body, err := ReqBase(url, reqData)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, errors.New("Json Unmarshal Error!!!")
	}
	if result.ErrorCode != 0 {
		return result, errors.New("Request Failed!!!")
	}
	return result, nil
}

func ReqBase(url string, reqData map[string]interface{}) (body []byte, err error) {
	var (
		client  *http.Client
		req     *http.Request
		resp    *http.Response
		newData []byte
	)
	client = &http.Client{}
	newData, _ = json.Marshal(reqData)
	req, err = http.NewRequest("POST", url, strings.NewReader(string(newData)))
	if err != nil {
		return nil, err
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
