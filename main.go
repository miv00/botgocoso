package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"os"
	"strings"
	"tgbotapi"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("no")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	pumpledb := initializeDb()

	updates := getUpdate(bot)

	for update := range updates {
		// mensaje no vacio recibido
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			chatID := update.Message.Chat.ID
			// si es start
			if update.Message.Text == "/start" {
				pumpleData, errRow := searchPumpleByID(pumpledb, chatID)
				if errRow != nil {
					cargarPumple(chatID, bot, pumpledb)
				}
				println(pumpleData.cumple)

				//Chequeo si es el cumple
				//Si cumple
				if pumpleData.cumple == time.Now().Format("02-04-2006") {
					felizcumple(chatID, bot, update.Message.Chat.UserName)
				}
				if pumpleData.cumple != time.Now().Format("02-04-2006") {
					felizNoCumple(chatID, bot)
				}
			}

			//REPLIES
			if update.Message.ReplyToMessage != nil {
				reply := strings.Split(update.Message.ReplyToMessage.Text, "para ")
				switch reply[0] {
				case "Hola ,me pasas tu cumple(DD-MM-AAAA)":
					anotarPumple(chatID, bot, update, pumpledb)
				case "Pasame el usuario":
					addGreetings(pumpledb, bot, update, chatID)
				case "Ingresar Saludo ":
					editGreeting(pumpledb, bot, update, chatID, reply[1])
				}

			}

		}
		//CALLBACKS QUERYS
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			commands := strings.Split(update.CallbackQuery.Data, ",")
			switch update.CallbackQuery.Data {

			case "bday":
				datos, _ := searchPumpleByID(pumpledb, chatID)
				fechaCumple := bdateParser(datos.cumple[0:10])
				bot.Send(howLongTillBday(fechaCumple, chatID))
			case "saludos":
				sendSaludosMenu(bot, chatID)
			case "cargar":
				requestUser(chatID, bot, update)
			case "borrar":
				callbackQuerylist(bot, pumpledb, "delete", chatID)
			case "cumple":
				cargarPumple(chatID, bot, pumpledb)
			case "editar":
				callbackQuerylist(bot, pumpledb, "update", chatID)
			case "falta":

			}
			if len(commands) > 1 {
				if commands[1] == "delete" {
					deleteGreetings(pumpledb, bot, commands[0], chatID)
				}
				if commands[1] == "update" {
					sendReplyUpdateGreetings(bot, commands[0], chatID)
				}
			}
		}

	}
}

func pumple(chat tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := chat.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Imprimiendo ... | "+" \n"+" ")
	edit, _ := bot.Send(msg)
	for _, line := range msj {
		var text string = edit.Text + "\n" + line
		coso := imprimiendo(text)
		edit, _ = bot.Send(tgbotapi.NewEditMessageText(chatID, edit.MessageID, coso))
	}
}

func sendSaludosMenu(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Menu Saludos")
	msg.ReplyMarkup = saludos
	bot.Send(msg)
}

func cargarPumple(chatID int64, bot *tgbotapi.BotAPI, pumpledb *sql.DB) {
	msg := tgbotapi.NewMessage(chatID, "Hola ,me pasas tu cumple(DD-MM-AAAA)")
	msg.ReplyMarkup = tgbotapi.ForceReply{true,
		"DD-MM-AAAA",
		false}
	bot.Send(msg)

}

func anotarPumple(chatID int64, bot *tgbotapi.BotAPI, update tgbotapi.Update, pumpledb *sql.DB) {
	bday, errDate := time.Parse("02-01-2006", update.Message.Text)
	if errDate != nil {
		msg := tgbotapi.NewMessage(chatID, "No entendi la fecha,Acordate que es DD-MM-AAA \n Si cumplo el 27 de agosto de 1983 \n Entonces ingreso 27-08-1983")
		msg.ReplyMarkup = anotarCumple
		bot.Send(msg)
	} else {
		ok := tgbotapi.NewMessage(chatID, "Pumple Anotado "+update.Message.Chat.UserName)
		bdayFinal := cumple{
			bday.String(),
			chatID,
			update.Message.Chat.UserName,
		}
		insertPumples(pumpledb, bdayFinal)
		bot.Send(ok)
	}
}

func bdateParser(cumple string) time.Time {
	fechaCumple, err := time.Parse("2006-01-02", cumple)
	println(err)
	fechaCumple = time.Date(time.Now().Year(), fechaCumple.Month(), fechaCumple.Day(), 0, 0, 0, 0, time.UTC)
	return fechaCumple
}

func imprimiendo(texto string) string {

	if strings.Contains(texto, "|") {
		texto = strings.Replace(texto, "|", "/", 1)
		texto = strings.Replace(texto, ".", " ", 1)
		return texto
	} else if strings.Contains(texto, "/") {
		texto = strings.Replace(texto, "/", "-", 1)
		texto = strings.Replace(texto, ".", " ", 1)
		return texto
	} else if strings.Contains(texto, "-") {
		texto = strings.Replace(texto, "-", "\\", 1)
		texto = strings.Replace(texto, ".", " ", 1)
		return texto
	} else if strings.Contains(texto, "\\") {
		texto = strings.Replace(texto, "\\", "|", 1)
		texto = strings.Replace(texto, " ", ".", 3)
		return texto
	}
	return texto
}

func initializeDb() *sql.DB {
	_, errdb := os.Stat("pumples.db")
	if os.IsNotExist(errdb) {
		CreateDbPumples()
	}
	//abro la base de datos
	pumpledb, err := sql.Open("sqlite3", "pumples.db")
	if err != nil {
		log.Panic(err)
	}
	return pumpledb
}

func getUpdate(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	return updates
}

func felizcumple(chatID int64, bot *tgbotapi.BotAPI, user string) {
	felicidades := tgbotapi.NewMessage(chatID, "Es tu pumple "+user)
	felicidades.ReplyMarkup = felizCumpleTeclado
	bot.Send(felicidades)
}

func felizNoCumple(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "Bot para desear feliz pumple")
	msg.ReplyMarkup = options
	bot.Send(msg)
}

func howLongTillBday(bday time.Time, chatID int64) tgbotapi.MessageConfig {
	if time.Now().After(bday) {
		bday = bday.AddDate(1, 0, 0)
	}
	duration := time.Until(bday)
	totalHours := duration.Hours()
	years := math.Trunc(totalHours / 8760)
	totalHours = totalHours - years*8760
	month := math.Trunc(duration.Hours() / 730)
	totalHours = totalHours - month*730
	weeks := math.Trunc(totalHours / 168)
	totalHours = totalHours - weeks*168
	days := math.Trunc(totalHours / 24)
	totalHours = totalHours - days*24
	msg := howLongTillBdayString(duration, years, month, weeks, days)
	return tgbotapi.NewMessage(chatID, msg)
}

func howLongTillBdayString(duration time.Duration, years float64, months float64, weeks float64, days float64) string {
	count := "Faltan "
	if years > 0 {
		count = count + fmt.Sprint(years) + "Año,"
	}
	if months > 0 {
		count = count + fmt.Sprint(months) + " Meses, "
	}
	if weeks > 0 {
		count = count + fmt.Sprint(weeks) + " Semanas, "
	}
	if days > 0 {
		count = count + fmt.Sprint(days) + " Dias y "
	}
	return count + duration.String()
}

func requestUser(chatID int64, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(chatID, "Pasame el usuario")
	msg.ReplyMarkup = tgbotapi.ForceReply{true,
		"DD-MM-AAAA",
		false}
	bot.Send(msg)

}

func getGreetings(chatID int64, bot *tgbotapi.BotAPI, pumpledb *sql.DB, text string) {
	_, err := searchPumpleByUser(pumpledb, text)
	if err != nil {
		notFound := tgbotapi.NewMessage(chatID, "No encontre al usuario  "+text)
		notFound.ReplyMarkup = errorAddGreetings
		bot.Send(notFound)
	} else {

	}

}

func addGreetings(pumpledb *sql.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update, chatID int64) {
	datos, err := searchPumpleByUser(pumpledb, update.Message.Text)
	if err != nil {
		notFound := tgbotapi.NewMessage(chatID, "No encontre al usuario  "+update.Message.Text)
		bot.Send(notFound)
	} else {
		_, errGreetings := alreadyAdded(pumpledb, datos.ChatID, chatID)
		if errGreetings != nil {
			greeting := saludo{
				saludo:           "Feliz Cumple de parte de ",
				Receiver:         update.Message.Text,
				Receiver_User_ID: datos.ChatID,
				Sender:           update.Message.Chat.UserName,
				Sender_User_ID:   chatID,
			}
			_, errInsert := insertSaludos(pumpledb, greeting)
			if errInsert != nil {
				print(errInsert)
			} else {
				succes := tgbotapi.NewMessage(chatID, "Saludo a "+update.Message.Text+" cargado correctamente")
				bot.Send(succes)
			}
		} else {
			alreadyExist := tgbotapi.NewMessage(chatID, "Ya tengo un saludo para esa persona")
			bot.Send(alreadyExist)
		}
	}
}

func callbackQuerylist(bot *tgbotapi.BotAPI, pumplesdb *sql.DB, data string, chatID int64) {
	greetings, err := searchLoadedGreetings(pumplesdb, chatID)
	if err != nil {
		log.Fatal(err)
	}
	var msg = tgbotapi.NewInlineKeyboardMarkup()
	if len(greetings) > 0 {
		for i := range greetings {
			button := tgbotapi.NewInlineKeyboardButtonData(greetings[i].Receiver, greetings[i].Receiver+","+data)
			row := tgbotapi.NewInlineKeyboardRow(button)
			msg.InlineKeyboard = append(msg.InlineKeyboard, row)
		}
		greatinList := tgbotapi.NewMessage(chatID, "Saludos Cargados")
		greatinList.ReplyMarkup = msg
		bot.Send(greatinList)
	} else {
		noGreetingss := tgbotapi.NewMessage(chatID, "Parece que no saludaste a nadie :(")
		bot.Send(noGreetingss)
	}
}

func deleteGreetings(pumple *sql.DB, bot *tgbotapi.BotAPI, data string, chatID int64) {
	_, err := deleteGreetingDB(pumple, data, chatID)
	if err != nil {
		log.Panic()
	} else {
		msg := tgbotapi.NewMessage(chatID, "Saludo a "+data+"borrado exitosamente")
		bot.Send(msg)
	}
}

func sendReplyUpdateGreetings(bot *tgbotapi.BotAPI, data string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Ingresar Saludo para "+data)
	msg.ReplyMarkup = tgbotapi.ForceReply{true,
		data,
		false}
	bot.Send(msg)
}

func editGreeting(pumpledb *sql.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update, chatID int64, receiver string) {
	replyMessage := update.Message.Text
	if replyMessage != "" {
		_, err := updateGreeting(pumpledb, replyMessage, receiver, chatID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatID, "El saludo estaba vacio")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Volver a intentar?", "Ingresar Saludo para "+receiver)))
		bot.Send(msg)
	}

}

var options = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Es Hoy?", "bday")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Saludos", "saludos")))

var saludos = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Cargar ", "Cargar")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Editar", "editar")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Borrar", "borrar")))

var cumpleTeclado = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Anotar Cumple", "cumple")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cambiar Fecha", "cambiar")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cuanto Falta?", "falta")))

var controls = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("<<", "list_<<"),
	tgbotapi.NewInlineKeyboardButtonData(">>", "list_>>"))

var felizCumpleTeclado = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Ver Saludos", "verSaludos")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ir al Otro Menu", "menu")))

var anotarCumple = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Anotar Cumple?", "cumple")))

var errorAddGreetings = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Volver a intetar?", "cargar")))

var msj = []string{
	"░▒█▀▀▀░█▀▀░█░░░▀░░▀▀█",
	"░▒█▀▀░░█▀▀░█░░░█▀░▄▀▒",
	"░▒█░░░░▀▀▀░▀▀░▀▀▀░▀▀▀",
	"▒█▀▄░█▒█░█▄▒▄█▒█▀▄░█▒░▒██▀",
	"░█▀▒░▀▄█░█▒▀▒█░█▀▒▒█▄▄░█▄▄",
	"༽΄◞ิ౪◟ิ‵༼ (｡◕‿‿◕｡)(づ｡◕‿‿◕｡)づ"}

type cumple struct {
	cumple string
	ChatID int64
	user   string
}

type saludo struct {
	saludo           string
	Receiver         string
	Receiver_User_ID int64
	Sender           string
	Sender_User_ID   int64
}

// ver variables go env
