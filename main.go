package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"strings"
	"github.com/joho/godotenv"
	"github.com/m1guelpf/chatgpt-telegram/src/chatgpt"
	"github.com/m1guelpf/chatgpt-telegram/src/config"
	"github.com/m1guelpf/chatgpt-telegram/src/session"
	"github.com/m1guelpf/chatgpt-telegram/src/tgbot"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Couldn't load config: %v", err)
	}

	if config.OpenAISession == "" {
		session, err := session.GetSession()
		if err != nil {
			log.Fatalf("Couldn't get OpenAI session: %v", err)
		}

		err = config.Set("OpenAISession", session)
		if err != nil {
			log.Fatalf("Couldn't save OpenAI session: %v", err)
		}
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load .env file: %v", err)
	}

	editInterval := 1 * time.Second
	if os.Getenv("EDIT_WAIT_SECONDS") != "" {
		editSecond, err := strconv.ParseInt(os.Getenv("EDIT_WAIT_SECONDS"), 10, 64)
		if err != nil {
			log.Printf("Couldn't convert your edit seconds setting into int: %v", err)
			editSecond = 1
		}
		editInterval = time.Duration(editSecond) * time.Second
	}

	bot, err := tgbot.New(os.Getenv("TELEGRAM_TOKEN"), editInterval)
	if err != nil {
		log.Fatalf("Couldn't start Telegram bot: %v", err)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bot.Stop()
		os.Exit(0)
	}()

	log.Printf("Started Telegram bot! Message @%s to start.", bot.Username)

	for update := range bot.GetUpdatesChan() {
		if update.Message == nil {
			continue
		}

		var (
//			updateText      = update.Message.Text
			updateChatID    = update.Message.Chat.ID
			updateMessageID = update.Message.MessageID
		)

		userId := strconv.FormatInt(update.Message.Chat.ID, 10)
		whiteLists := strings.Split(os.Getenv("TELEGRAM_ID"), ",")
		log.Printf("[Chat ID] %s", userId)
		if !(len(whiteLists) == 1 && whiteLists[0] == "") && !contains(whiteLists, userId) {
			bot.Send(updateChatID, updateMessageID, "No estás autorizado para usar este bot | You are not authorized to use this bot.")
			continue
		}

		if !update.Message.IsCommand(){
			if (update.Message.Chat.IsPrivate() || ((update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) && strings.HasPrefix(update.Message.Text, "@" + bot.Username))) {
				file,err := bot.GetFileDirectURL(update.Message.Voice.FileID)
				if err != nil{
					fmt.Sprintf("Error: %v", err)
				}else{
					bot.SendTyping(updateChatID)
					cmd := exec.Command("python3", "./whisperAudio.py", file)
					stdout,err := cmd.Output()
					if err != nil {
						fmt.Println(err.Error())
						continue
    					}
					bot.Send(updateChatID, updateMessageID, string(stdout))
				}
				continue
			}else{
				continue
			}
		}

		var text string
		switch update.Message.Command() {
		case "help":
			text = "Envía un mensaje para empezar a hablar con ChatGPT. Puedes usar /reload en cualquier momento para empezar una nueva conversación desde cero (no se borrarán los mensajes de Telegram)."
		case "start":
			text = "Envía un mensaje para empezar a hablar con ChatGPT. Puedes usar /reload en cualquier momento para empezar una nueva conversación desde cero (no se borrarán los mensajes de Telegram)."
		case "reload":
			text = "Nueva conversación empezada. ¡Disfruta!"
		default:
			text = "Comando desconocido. Envía /help para ver una lista de comandos."
		}

		if _, err := bot.Send(updateChatID, updateMessageID, text); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
