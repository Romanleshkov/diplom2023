package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func sendMessage(botUrl string, botMessage *BotMessage) (*http.Response, error){
	buf, err := json.Marshal(*botMessage)
	if err != nil{
		return nil, err
	}
	resp, err := http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil{
		return nil, err
	}
	return resp, err
}

func editMessageText(botUrl string, botMessage *BotMessage) (*http.Response, error){
	buf, err := json.Marshal(botMessage)
	if err != nil{
		return nil, err
	}
	resp, err := http.Post(botUrl + "/editMessageText", "application/json",
		bytes.NewBuffer(buf))
	if err != nil{
		return nil, err
	}
	return resp, err
}

func answerCallbackQuery(botUrl string, answerQuery *AnswerCallbackQuery) (*http.Response, error){
	buf, err := json.Marshal(answerQuery)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(botUrl + "/answerCallbackQuery", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return resp, err
}