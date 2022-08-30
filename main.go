package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"tgbotapi"
	"time"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("quemiras")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, errdb := os.Stat("pumples.db")
	if os.IsNotExist(errdb) {
		CreateDbPumples()
	}
	//abro la base de datos
	pumpledb, err := sql.Open("sqlite3", "pumples.db")
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// mensaje no vacio recibido
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			chatID := update.Message.Chat.ID
			// si es start
			if update.Message.Text == "/start" {
				pumpleData, errRow := searchPumpleByID(pumpledb, chatID)

				if errRow != nil {
					// cumpleData := cumple{"asdas", update.Message.Chat.ID, update.Message.Chat.UserName}
					cargarPumple(chatID, bot, pumpledb)
				}
				println(pumpleData.cumple)

				//Chequeo si es el cumple
				if pumpleData.cumple == time.Now().Format("02-04-2006") {
					felicidades := tgbotapi.NewMessage(chatID, "Es tu pumple "+update.Message.Chat.UserName)
					felicidades.ReplyMarkup = felizCumpleTeclado
					bot.Send(felicidades)

				}
				if pumpleData.cumple != time.Now().Format("02-04-2006") {
					msg := tgbotapi.NewMessage(chatID, "Bot para desear feliz pumple")
					msg.ReplyMarkup = options
					bot.Send(msg)
				}
			}

			//REPLIES
			if update.Message.ReplyToMessage != nil {
				switch update.Message.ReplyToMessage.Text {
				case "Hola ,me pasas tu cumple(DD-MM-AAAA)":
					bday, errDate := time.Parse("02-01-2006", update.Message.Text)
					if errDate != nil {
						msg := tgbotapi.NewMessage(chatID, "No entendi la fecha,Acordate que es DD-MM-AAA \n Si cumplo el 27 de agosto de 1983 \n Entonces ingreso 27-08-1983")
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Anotar Cumple?", "cumple")))
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
				case "Pasame el usuario":
					_, err = searchPumpleByUser(pumpledb, update.Message.Text)
					if err != nil {
						notFound := tgbotapi.NewMessage(chatID, "No encontre al usuario  "+update.Message.Text)
						bot.Send(notFound)
					} else {

					}

				}

			}

		}
		//CALLBACKS QUERYS
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			switch update.CallbackQuery.Data {

			case "bday":
				datos, _ := searchPumpleByID(pumpledb, chatID)
				fechaCumple, err := time.Parse("2006-01-02", datos.cumple[0:10])
				println(err)
				fechaCumple = time.Date(time.Now().Year(), fechaCumple.Month(), int(fechaCumple.Day()), 0, 0, 0, 0, time.UTC)
				if time.Now().Before(fechaCumple) {
					until := time.Until(fechaCumple)
					days := "Faltan " + until.String()
					time := tgbotapi.NewMessage(chatID, days)
					bot.Send(time)
				} else {
					fechaCumple = fechaCumple.AddDate(1, 0, 0)
					until := time.Until(fechaCumple)
					days := "Faltan " + until.String()
					time := tgbotapi.NewMessage(chatID, days)
					bot.Send(time)
				}
				msg := tgbotapi.NewMessage(chatID, "...Hi")
				bot.Send(msg)

			case "saludo":
				msg := tgbotapi.NewMessage(chatID, "Pasame el usuario")
				msg.ReplyMarkup = tgbotapi.ForceReply{true,
					"DD-MM-AAAA",
					false}
				bot.Send(msg)
			case "borrar":

			case "cumple":
				cargarPumple(chatID, bot, pumpledb)
			case "cambiar":
				msg := tgbotapi.NewMessage(chatID, "este si anda)")
				bot.Send(msg)
			case "falta":

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

func cargarPumple(chatID int64, bot *tgbotapi.BotAPI, pumpledb *sql.DB) {
	msg := tgbotapi.NewMessage(chatID, "Hola ,me pasas tu cumple(DD-MM-AAAA)")
	msg.ReplyMarkup = tgbotapi.ForceReply{true,
		"DD-MM-AAAA",
		false}
	bot.Send(msg)

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

var options = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Es Hoy?", "bday")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cargar Saludo", "saludo")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Borrar Saludo", "borrar")))

var cumpleTeclado = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Anotar Cumple", "cumple")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cambiar Fecha", "cambiar")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Cuanto Falta?", "falta")))

var felizCumpleTeclado = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Ver Saludos", "verSaludos")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ir al Otro Menu", "menu")))

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
	saludo  string
	ChatID  int64
	User_ID string
}

// ver variables go env
