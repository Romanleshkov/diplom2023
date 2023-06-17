package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	Id int				`json:"id"`
	Login string		`json:"login"`
	Password string		`json:"password"`
	Server string		`json:"server"`
}

func checkCacheDir() error{
	_, err := os.Stat("./cache")
	if err == nil { return nil }
	if os.IsNotExist(err){
		err = os.MkdirAll("./.cache", 0700)
		if err != nil{
			return err
		}
		return nil
	}
	return err
}

func checkUserData(userId int) (bool,error){
	if err := checkCacheDir(); err != nil{
		log.Println("Error in func checkCacheDir(): ", err)
	}
	userPath := "./.cache/" + strconv.Itoa(userId)
	_, err := os.Stat(userPath)
	if err == nil{
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist){
		return false, nil
	}
	return false, err
}

func getUserData(userId int) (UserData,error) {
	userPath := "./.cache/" + strconv.Itoa(userId)
	if fileExist, err := checkUserData(userId); err != nil{
		log.Println("Error in func checkUserPath(): ", err)
		return UserData{}, err
	}else if fileExist{
		file, err := os.Open(userPath)
		if err != nil{
			return UserData{}, err
		}
		defer func() {
			if err := file.Close(); err != nil{
				log.Println("Error from func file.Close(): ", err)
			}
		}()
		body, err := ioutil.ReadAll(file)

		var userData UserData
		if err := json.Unmarshal(body, &userData); err != nil{
			return UserData{}, err
		}
		return userData, nil
	}else{
		return UserData{}, nil
	}
}

func createUserData(userId int, data UserData) error{
	userPath := "./.cache/" + strconv.Itoa(userId)
	if _, err := checkUserData(userId); err != nil{
		log.Println("Error in func checkUserPath(): ", err)
		return err
	}else{
		file, err := os.Create(userPath)
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
			return err
		}
		defer func() {
			if err := file.Close(); err != nil{
				log.Println("Error from func file.Close(): ", err)
			}
		}()
		dataJson, err := json.MarshalIndent(data, "", "\t")
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
		}
		_, err = file.Write(dataJson)
		if err != nil{
			log.Println("Error in func createUserData(): ", err)
		}
		return nil
	}
}

func strToUserData(str string) (UserData,error){
	str2 := strings.Split(str, "\n")
	if len(str2) != 4{
		return UserData{}, errors.New("wrong format of userData")
	}
	var flag = 0
	var userData UserData

	fmt.Println(str2)
	for i, l := range str2{
		switch i {
		case 0:
			if !strings.HasPrefix(l, "id: "){
				return UserData{}, errors.New("wrong format of userData")
			}
			id, err := strconv.Atoi(l[strings.Index(l," ") + 1:])
			if err != nil{
				return UserData{}, errors.New("wrong format of userData")
			}
			userData.Id = id
		case 1:
			if !strings.HasPrefix(l, "login: "){
				return UserData{}, errors.New("wrong format of userData")
			}
			userData.Login = l[strings.Index(l," ") + 1:]
		case 2:
			if !strings.HasPrefix(l, "password: "){
				return UserData{}, errors.New("wrong format of userData")
			}
			userData.Password = l[strings.Index(l," ") + 1:]
		case 3:
			if !strings.HasPrefix(l, "server: "){
				return UserData{}, errors.New("wrong format of userData")
			}
			userData.Server = l[strings.Index(l," ") + 1:]
		}
		flag ++
		}
	return userData, nil
}

func (userData UserData)  String() string{
	return "\nТекущие данные:\n" +
		"id: " + strconv.Itoa(userData.Id) + "\n" +
		"login: " + userData.Login + "\n" +
		"password: " + userData.Password + "\n" +
		"server: " + userData.Server + "\n"
}

func strUnixTimeToStrDate(timeUnix string) string {
	tm,err := strconv.ParseInt(timeUnix, 10, 64)
	if err != nil{
		log.Println("Error in func strconv(): ", err)
		return "Error"
	}
	tmTime := time.Unix(tm,0)
	tmStr := tmTime.Format("02.01.2006")
	return tmStr
}

