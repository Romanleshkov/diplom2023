package main

import (
	"bytes"
	//"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	//"strings"
)

func getAuthToken(id int, username string, password string, server string) (UserLoginResponse, error) {
	userLoginRequest := newUserLoginRequest(id, username, password)
	zabbixURL := "http://" + server + "/api_jsonrpc.php"
	buf, err := json.Marshal(userLoginRequest)
	if err != nil{
		log.Println("Error in func json.Marshal(): ", err)
		return UserLoginResponse{}, err
	}
	resp, err := http.Post(zabbixURL, "application/json", bytes.NewBuffer(buf))
	if err != nil{
		log.Println("Error in func http.Post(): ", err)
		return UserLoginResponse{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil{
			log.Println("Error from func resp.Body.Close(): ", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("Error in func ioutil.ReadAll(): ", err)
		return UserLoginResponse{}, err
	}
	var response UserLoginResponse
	if err := json.Unmarshal(body, &response); err != nil{
		log.Println("Error in func json.Unmarshal(): ", err)
		return UserLoginResponse{}, err
	}
	if response.Error.Message != ""{
		return UserLoginResponse{}, errors.New(response.Error.Data)
	}
	return response, nil
}

func ProblemGet(TGUserId int) (EventGetResponse, error){
	var request = newProblemGetRequest()
	var response EventGetResponse
	err := zabbixAuthRequest(TGUserId, &request, &response)
	if err != nil{
		log.Println("Error in func zabbixAuthRequest(): ", err)
		return EventGetResponse{}, err
	}
	return response, nil
}

func eventGet(TGUserId int, eventId string) (EventGetResponse, error){
	var request = newEventGetRequest()
	var response EventGetResponse
	request.Params.Eventids = append(request.Params.Eventids, eventId)
	err := zabbixAuthRequest(TGUserId, &request, &response)
	if err != nil{
		log.Println("Error in func zabbixAuthRequest(): ", err)
		return EventGetResponse{}, err
	}
	return response, nil
}

func eventAcknowledge(TGUserId int, eventId string, actions int, severity int, message string) error{
	var request = newEventAcknowledgeRequest()
	var response = struct {
		Error ErrorResponse	`json:"error"`
	}{}
	request.Params.Eventids = eventId
	request.Params.Action = actions
	request.Params.Severity = severity
	request.Params.Message = message
	err := zabbixAuthRequest(TGUserId, &request, &response)
	if err != nil{
		log.Println("Error in func zabbixAuthRequest(): ", err)
		return err
	}
	if response.Error.Message != ""{
		log.Println("Error response from Zabbix in eventAcknowledge(): ", response)
		return err
	}
	return nil
}

func scriptGetscriptsbyevents(TGUserId int, eventId string) (ScriptGetscriptsbyeventsResponse, error){
	var request = newScriptGetscriptsbyeventsRequest()
	var response ScriptGetscriptsbyeventsResponse
	request.Params.Eventids = eventId
	err := zabbixAuthRequest(TGUserId, &request, &response)
	if err != nil{
		log.Println("Error in func zabbixAuthRequest(): ", err)
		return ScriptGetscriptsbyeventsResponse{}, err
	}
	return response, nil
}

func scriptExecute(TGUserId int, scriptId string, eventId string) (ScriptExecuteResponse, error){
	var request = newScriptExecuteRequest()
	var response ScriptExecuteResponse
	request.Params.Eventid = eventId
	request.Params.Scriptid = scriptId
	err := zabbixAuthRequest(TGUserId, &request, &response)
	if err != nil{
		log.Println("Error in func zabbixAuthRequest(): ", err)
		return ScriptExecuteResponse{}, err
	}
	return response, nil
}

func zabbixAuthRequest(TGUserId int, request interface{setZabbixId(int)}, response interface{}) error{
	userData, err := getUserData(TGUserId)
	if err != nil{
		log.Println("Error in func getUserData()")
	}
	authToken, err := getAuthToken(userData.Id, userData.Login, userData.Password, userData.Server)
	if err != nil {
		return err
	}
	zabbixURL := "http://" + userData.Server + "/api_jsonrpc.php"

	request.setZabbixId(userData.Id)
	buf, err := json.Marshal(request)
	if err != nil{
		log.Println("Error in func json.Marshal(): ", err)
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", zabbixURL, bytes.NewBuffer(buf))
	if err != nil{
		log.Println("Error in func http.NewRequest(): ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + authToken.Result)
	resp, err := client.Do(req)
	if err != nil{
		log.Println("Error in func client.Do(): ", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil{
			log.Println("Error from func resp.Body.Close(): ", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println("Error in func ioutil.ReadAll(): ", err)
		return err
	}
	if err := json.Unmarshal(body, &response); err != nil{
		log.Println("Error in func json.Unmarshal(): ", err)
		return err
	}
	return nil
}
