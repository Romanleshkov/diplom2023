package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main(){
	if err := godotenv.Load("config.env"); err != nil{
		log.Fatal("file is not found: config.env")
	}
	// https://api.telegram.org/bot<token>/METHOD_NAME
	botToken := os.Getenv("BOT_TOKEN")
	botUrl := "https://api.telegram.org/bot" + botToken
	offset := 0
	for true{
		updates, err := getUpdates(botUrl, offset)
		if err != nil{
			log.Println("Error from func getUpdates: ", err.Error())
		}
		for _, update := range updates{
			fmt.Println("update: ", update)
			go func() {
				if err := respondChat(botUrl, update); err != nil{
					log.Println("Error from func respondChat: ", err.Error())
				}
			}()
			go func() {
				if err := respondButton(botUrl, update); err != nil{
					log.Println("Error from func respondQuery: ", err.Error())
				}
			}()
			offset = update.UpdateId + 1
		}
	}

}

func getUpdates(botUrl string, offset int) ([]Update, error){
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil{
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil{
			log.Println("Error from func resp.Body.Close(): ", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var response Response
	if err := json.Unmarshal(body, &response); err != nil{
		return nil, err
	}
	return response.Result, nil
}

func respondChat(botUrl string, update Update) error {
	if update.Message.Chat.Id == 0{
		return nil
	}
	var botMessage BotMessage
	err := commandChoice(update.Message.Chat.Id, update.Message.Text, &botMessage)
	if err != nil {
		log.Println("Error in func commandChoice(): ", err)
	}
	_, err = sendMessage(botUrl, &botMessage)
	if err != nil{
		return err
	}
	return nil
}

func respondButton (botUrl string, update Update) error {
	if update.CallbackQuery.Message.Chat.Id == 0 {
		return nil
	}
	var answerQuery AnswerCallbackQuery
	answerQuery.CallbackQueryId = update.CallbackQuery.Id
	answerQuery.Text = ""
	var editMessage BotMessage
	editMessage.MessageId = update.CallbackQuery.Message.MessageId
	err := commandChoice(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Data, &editMessage)
	if err != nil {
		log.Println("Error in func commandChoice(): ", err)
	}
	_, err = editMessageText(botUrl, &editMessage)
	if err != nil{
		return err
	}
	_, err = answerCallbackQuery(botUrl, &answerQuery)
	if err != nil{
		return err
	}
	return nil
}

func commandChoice(TGChatId int, text string, message *BotMessage) error{
	message.ChatId = TGChatId
	out := mainMenu(TGChatId)
	message.Text = out.text
	message.ReplyMarkup.InlineKeyboard = out.keyboard

	command := strings.Split(text, " ")
	switch command[0] {
	case "/start":

	case "/mainMenu":

	case "/userDataMenu":
		out := userDataMenu(TGChatId)
		message.Text = out.text
		message.ReplyMarkup.InlineKeyboard = out.keyboard
	case "/user":
		if len(command) == 6 {
			userData, err := strToUserData(strings.Join(command[1:], " "))
			if err != nil {
				out := userDataMenu(TGChatId)
				message.Text = "Неверный формат пользовательских данных\n" + out.text
				message.ReplyMarkup.InlineKeyboard = out.keyboard
			} else {
				out := mainMenu(TGChatId)
				if err := createUserData(TGChatId, userData); err != nil {
					message.Text = out.text + "\nОшибка при создании пользователя.\n"
					log.Println("Error in func createUserData(): ", err)
				} else {
					message.Text = out.text + "\nДанные внесены.\n" + userData.String()
				}
				message.ReplyMarkup.InlineKeyboard = out.keyboard
			}
		}else {
			message.Text += "\nНедостаточно параметров для команды /user"
		}
	case "/problem.get":
		resp, err := ProblemGet(TGChatId)
		if err != nil {
			log.Println("Error in func problemGet(): ", err)
		}
		out := problemGetMenu(resp)
		message.Text = out.text
		message.ReplyMarkup.InlineKeyboard = out.keyboard
	case "/event.get":
		if len(command) == 2 {
			resp, err := eventGet(TGChatId, command[1])
			if err != nil {
				log.Println("Error in func eventGet(): ", err)
			}
			out := eventGetMenu(resp)
			message.Text = out.text
			message.ReplyMarkup.InlineKeyboard = out.keyboard
		}

	case "/event.acknowledge":
		if len(command) == 2 {
			out := eventAcknowledgeMenu(command[1], 0)
			message.Text = out.text
			message.ReplyMarkup.InlineKeyboard = out.keyboard
		}
		if len(command) == 3 {
			actions, _ := strconv.Atoi(command[2])
			out := eventAcknowledgeMenu(command[1], actions)
			message.Text = out.text
			message.ReplyMarkup.InlineKeyboard = out.keyboard
		}
	case "/eventAcknowledgeInit":
		if len(command) == 3 {
			resp, err := eventGet(TGChatId, command[1])
			if err != nil {
				log.Println("Error in func eventGet(): ", err)
			} else {
				severity, _ := strconv.Atoi(resp.Result[0].Severity)
				actions, _ := strconv.Atoi(command[2])
				out := eventAcknowledgeInitMenu(command[1], actions, severity)
				message.Text = out.text
				message.ReplyMarkup.InlineKeyboard = out.keyboard
			}
		}
		if len(command) == 4 {
			severity, _ := strconv.Atoi(command[3])
			actions, _ := strconv.Atoi(command[2])
			out := eventAcknowledgeInitMenu(command[1], actions, severity)
			message.Text = out.text
			message.ReplyMarkup.InlineKeyboard = out.keyboard
		}
	case "/eventAcknowledgeFinish":
		if len(command) == 4 {
			actions, _ := strconv.Atoi(command[2])
			severity, _ := strconv.Atoi(command[3])
			err := eventAcknowledge(TGChatId, command[1], actions, severity, "")
			if err != nil {
				log.Println("Error in func eventAcknowledge(): ", err)
				message.Text += "\nError\n"
			} else {
				message.Text += "\nOK\n"
			}
		}
		if len(command) == 5{
			actions, _ := strconv.Atoi(command[2])
			severity, _ := strconv.Atoi(command[3])
			err := eventAcknowledge(TGChatId, command[1], actions, severity, command[4])
			if err != nil {
				log.Println("Error in func eventAcknowledge(): ", err)
				message.Text += "\nError\n"
			} else {
				message.Text += "\nOK\n"
			}
		}
	case "/scriptGetscriptsbuevents":
		if len(command) == 2 {
			resp, err := scriptGetscriptsbyevents(TGChatId, command[1])
			if err != nil {
				log.Println("Error in func eventGet(): ", err)
			} else {
				out := scriptGetscriptsbueventsMenu(resp, command[1])
				message.Text = out.text
				message.ReplyMarkup.InlineKeyboard = out.keyboard
			}
		}
	case "/scriptExecute":
		if len(command) == 3{
			resp, err := scriptExecute(TGChatId, command[1], command[2])
			if err != nil {
				log.Println("Error in func eventGet(): ", err)
			} else {
				out := scriptExecuteMenu(resp)
				message.Text = out.text
				message.ReplyMarkup.InlineKeyboard = out.keyboard
			}
		}
	default:
		message.Text += "Неизвестная команда"
	}
	return nil
}
