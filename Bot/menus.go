package main

import (
	"log"
	"strconv"
)

type menu struct {
	text string
	keyboard [][]InlineKeyboardButton
}

func userDataMenu(id int) menu{
	var menu = menu{
		text: "Укажите id, логин, пароль и сервер\n" +
			"В формате:\n" +
			"/user id: <id>\n" +
			"login: <login>\n" +
			"password: <password>\n" +
			"server: <serverIP:port or servername>\n",
		keyboard: [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{Text: "Главное меню", CallbackData: "/mainMenu"},
			},
		},
	}
	if fileExist, err := checkUserData(id); err != nil {
		log.Println("Error in func checkUserData(): ", err)
	} else if fileExist {
		userData, err := getUserData(id)
		if err != nil {
			log.Println("Error in func getUserData(): ", err)
		}
		menu.text += userData.String()
	}
	return menu
}

func mainMenu (id int) menu {
	var menu = menu{
		text: "Главное меню\n",
	}
	if ans, _ := checkUserData(id); ans {
		menu.text += "Используйте эту команду для обработки события:\n" +
			"\"/event.get <Original problem ID>\"\n" +
			"Или выберите одну из текущих проблем"
		menu.keyboard = [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{Text: "Изменить данные",
					CallbackData: "/userDataMenu"},
				InlineKeyboardButton{Text: "Текущие проблемы",
					CallbackData: "/problem.get"},
			},
		}
		return menu
	}
	menu.keyboard = [][]InlineKeyboardButton{
		{
			InlineKeyboardButton{Text: "Изменить данные",
				CallbackData: "/userDataMenu"},
		},
	}
	return menu
}

func problemGetMenu(problems EventGetResponse) menu {
	menu := menu{
		text: "Текущие проблемы\n",
			}
	for _, problem := range problems.Result{
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				{Text: problem.Name,
					CallbackData: "/event.get " + problem.Eventid}},
			)
	}
	if len(problems.Result) == 0{
		menu.text = "Проблем нет"
	}
	if problems.Jsonrpc == ""{
		menu.text = "Сервер недоступен"
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			{Text: "Главное меню",
				CallbackData: "/mainMenu"}},
	)
	return menu
}

func eventGetMenu(events EventGetResponse) menu{
	event := events.Result[0]
	tmStr := strUnixTimeToStrDate(event.Clock)
	menu := menu{
		text: "Информация о событии " + event.Eventid + "\n" +
			"Имя события: " + event.Name + "\n" +
			"Время обнаружения: " + tmStr + "\n" +
			"Важность: " + event.Severity + "\n",
	}
	if len(event.Acknowledges) > 0{
		menu.text += "Обновления события:\n"
	}
	for _, acknowledge := range event.Acknowledges{
		menu.text += "---:\n" +
			"Время: " + strUnixTimeToStrDate(acknowledge.Clock) + "\n" +
			"От: " + acknowledge.Username + "\n" +
			"Сообщение:\n" + acknowledge.Message + "\n"
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			{Text: "Скрипты",
				CallbackData: "/scriptGetscriptsbuevents " + event.Eventid},
			{Text: "Действия",
				CallbackData: "/event.acknowledge " + event.Eventid},
			{Text: "Главное меню",
				CallbackData: "/mainMenu"},
		},
	)
	return menu
}

func eventAcknowledgeMenu(eventId string, actions int) menu{
	menu := menu{
		text: "Выберите действие для " + eventId + "\n" +
			"После выбора нажмите \"Параметры\"\n",
	}
	var button1 = InlineKeyboardButton{Text: "Close problem: No",
			CallbackData: "/event.acknowledge " + eventId + " " + strconv.Itoa(1 ^ actions)}
	var button2 = InlineKeyboardButton{Text: "Add message: No",
			CallbackData: "/event.acknowledge " + eventId + " " + strconv.Itoa(4 ^ actions)}
	var button3 = InlineKeyboardButton{Text: "Change severity: No",
			CallbackData: "/event.acknowledge " + eventId + " " + strconv.Itoa(8 ^ actions)}
	if (1 ^ actions) < actions {
		button1.Text = "Close problem: Yes"
	}
	if (4 ^ actions) < actions {
		button2.Text = "Add message: Yes"
	}
	if (8 ^ actions) < actions {
		button3.Text = "Change severity: Yes"
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			button1,
			button2,
		},
		[]InlineKeyboardButton{
			button3,
		},
	)
	if actions == 0{
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				{Text: "Главное меню",
					CallbackData: "/mainMenu"},
			},
		)
	}else {
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				{Text: "Параметры",
					CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions)},
				{Text: "Главное меню",
					CallbackData: "/mainMenu"},
			},
		)
	}
	return menu
}

func eventAcknowledgeInitMenu(eventId string, actions int, severity int) menu {
	menu := menu{
		text: "Событие  " + eventId + "\n",
	}

	if (8 ^ actions) < actions {
		var button0 = InlineKeyboardButton{Text: "not classified: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 0"}
		var button1 = InlineKeyboardButton{Text: "information: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 1"}
		var button2 = InlineKeyboardButton{Text: "warning: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 2"}
		var button3 = InlineKeyboardButton{Text: "average: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 3"}
		var button4 = InlineKeyboardButton{Text: "high: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 4"}
		var button5 = InlineKeyboardButton{Text: "disaster: No",
			CallbackData: "/eventAcknowledgeInit " + eventId + " " + strconv.Itoa(actions) + " 5"}
		switch severity {
		case 0:
			button0.Text = "not classified: Yes"
		case 1:
			button1.Text = "information: Yes"
		case 2:
			button2.Text = "warning: Yes"
		case 3:
			button3.Text = "average: Yes"
		case 4:
			button4.Text = "high: Yes"
		case 5:
			button5.Text = "disaster: Yes"
		}
		menu.text += "\nВыбирите важность события\n"
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				button0,button1,button2,
			},
			[]InlineKeyboardButton{
				button3,button4,button5,
			},
		)
	}
	if (4 ^ actions) < actions{
		menu.text += "Для отправки сообщения и обновления события введите в строку следующую комманду:\n" +
			"\"/eventAcknowledgeFinish " + eventId + " " + strconv.Itoa(actions) + " " + strconv.Itoa(severity) + " <your message>\""
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				{Text: "Главное меню",
					CallbackData: "/mainMenu"},
			},
		)
	} else {
		menu.keyboard = append(
			menu.keyboard,
			[]InlineKeyboardButton{
				{Text: "Обновить",
					CallbackData: "/eventAcknowledgeFinish " + eventId + " " + strconv.Itoa(actions) + " " + strconv.Itoa(severity)},
				{Text: "Главное меню",
					CallbackData: "/mainMenu"},
			},
		)
	}
	return menu
}

func scriptGetscriptsbueventsMenu(scripts ScriptGetscriptsbyeventsResponse, eventId string) menu{
	menu := menu{
		text: "Выберите скрипт для события " + eventId + "\n",
	}
	for _, event := range scripts.Result{
		for _, script := range event {
			menu.keyboard = append(
				menu.keyboard,
				[]InlineKeyboardButton{
					InlineKeyboardButton{Text: "Выполнить скрипт: " + script.Name,
						CallbackData: "/scriptExecute " + script.Scriptid + " " + eventId}},
			)
		}
	}
	menu.keyboard = append(
		menu.keyboard,
		[]InlineKeyboardButton{
			{Text: "Главное меню",
				CallbackData: "/mainMenu"}},
	)
	return menu
}

func scriptExecuteMenu(ans ScriptExecuteResponse) menu{
	status := ans.Result.Response
	data := ans.Result.Value
	if ans.Error.Message != ""{
		status = "error"
		data = ans.Error.Data
	}
	menu := menu{
		text: "Статус выполнения скрипта: " + status + "\n" +
			"Вывод ответа скрипта:\n\n" +
			data,
		keyboard: [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{Text: "Главное меню", CallbackData: "/mainMenu"},
			},
		},
	}
	return menu
}
