package core

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MgdContext struct {
	// 签到相关信息
	Info CheckInfo
	// 登录后的Token
	Token string
	// planId
	planId string
}

// 请求的简单封装
func (m MgdContext) request(api string, requestData map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	httpReq, _ := http.NewRequest("POST", api, strings.NewReader(string(data)))
	httpReq.Header.Add("Content-Type", "application/json")
	if m.Token != "" {
		httpReq.Header.Add("Authorization", m.Token)
		// log.Println(m.Token)
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mapResult map[string]interface{}
	json.Unmarshal(b, &mapResult)

	// log.Println(mapResult)
	if int32(mapResult["code"].(float64)) != 200 {
		return nil, errors.New(fmt.Sprintf("接口返回异常 %v", mapResult["msg"]))
	}

	return mapResult, nil
}

// 登录
func (m *MgdContext) Login() error {
	api := "https://api.moguding.net:9000/session/user/v1/login"
	param := map[string]interface{}{
		"loginType": "android",
		"uuid":      uuid.NewV4().String(),
		"phone":     m.Info.Phone,
		"password":  m.Info.Password,
	}

	respMap, err := m.request(api, param)
	if err != nil {
		return err
	}

	token := (respMap["data"].(map[string]interface{}))["token"].(string)
	m.Token = token

	err = m.getPlanId()
	if err != nil {
		return err
	}

	return nil
}

// 考勤
func (m MgdContext) attendance(paramType string) error {
	api := "https://api.moguding.net:9000/attendence/clock/v1/save"
	param := map[string]interface{}{
		"type":        paramType,
		"device":      "Android",
		"country":     m.Info.Country,
		"province":    m.Info.Province,
		"city":        m.Info.City,
		"address":     m.Info.Address,
		"latitude":    m.Info.Latitude,
		"longitude":   m.Info.Longitude,
		"description": m.Info.Description,
		"planId":      m.planId,
	}
	respMap, err := m.request(api, param)
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(respMap)

	if int32(respMap["code"].(float64)) != 200 {
		return errors.New(fmt.Sprintf("接口返回异常 %v", respMap["msg"]))
	}
	return nil
}

// 获取 PlanId
func (m *MgdContext) getPlanId() error {
	api := "https://api.moguding.net:9000/practice/plan/v1/getPlanByStu"
	param := map[string]interface{}{
		"statue": "",
	}
	respMap, err := m.request(api, param)
	if err != nil {
		return err
	}

	m.planId = respMap["data"].([]interface{})[0].(map[string]interface{})["planId"].(string)
	return nil
}

// 上班打卡
func (m MgdContext) StartWork() error {
	err := m.attendance("START")
	if err != nil {
		return err
	}
	return nil
}

// 下班打卡
func (m MgdContext) EndWork() error {
	err := m.attendance("END")
	if err != nil {
		return err
	}
	return nil
}
