package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot *tgbotapi.BotAPI

	homeReplyKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Tell me who you are"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Working career"),
			tgbotapi.NewKeyboardButton("Technologies and Projects"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Contacts"),
		),
	)
)

type Messages map[string][]*Message

// Message - export
type Message struct {
	ChatID   int64                         `json:"chat_id" yaml:"chat-id"`
	MsgType  string                        `json:"msg_type" yaml:"msg-type"`
	Duration time.Duration                 `json:"duration" yaml:"duration"`
	Content  string                        `json:"content" yaml:"content"`
	Keyboard *tgbotapi.ReplyKeyboardMarkup `json:"-" yaml:"-"`
}

func consumeChainMessage(structure Message) {
	switch structure.MsgType {
	case "Message":
		var response tgbotapi.MessageConfig = tgbotapi.NewMessage(structure.ChatID, structure.Content)
		response.ParseMode = "Markdown"
		if structure.Keyboard != nil {
			response.ReplyMarkup = structure.Keyboard
		}

		SendMsg(response)

	case "NewDocumentUpload":
		var response tgbotapi.DocumentConfig = tgbotapi.NewDocumentUpload(structure.ChatID, structure.Content)
		if structure.Keyboard != nil {
			response.ReplyMarkup = structure.Keyboard
		}

		SendMsg(response)

	case "NewPhotoUpload":
		var response tgbotapi.PhotoConfig = tgbotapi.NewPhotoUpload(structure.ChatID, structure.Content)
		if structure.Keyboard != nil {
			response.ReplyMarkup = structure.Keyboard
		}

		SendMsg(response)
	}

	time.Sleep(structure.Duration * time.Second)

}

func updatesHandler() {
	var telegramApikey string

	errDotEnv := godotenv.Load()
	if errDotEnv != nil {
		log.Fatal("Error loading .env file")
	}
	telegramApikey = os.Getenv("TELEGRAM_APIKEY")

	var err error
	bot, err = tgbotapi.NewBotAPI(telegramApikey)
	bot.Debug = true

	if err != nil {
		log.Panicln(err)
	}

	logMessage := fmt.Sprintf("Bot connected correctly %s", bot.Self.UserName)
	log.Println(logMessage)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, chanErr := bot.GetUpdatesChan(u)
	if chanErr != nil {
		log.Panicln(chanErr)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		commandsHandler(update)
	}
}

func commandsHandler(update tgbotapi.Update) {
	command, _, ok := breakCommand(update.Message.Text)
	if ok {
		switch command {
		case "Back":
			TornaCommand(update)
		case "Start":
			StartCommand(update)
		case "/Start":
			StartCommand(update)
		case "/start":
			StartCommand(update)
		case "tellme":
			StoryCommand(update)
		case "Track":
			JobsCommand(update)
		case "Tecnologies":
			TechCommand(update)
		case "Contats":
			ContactsCommand(update)
		}
	}
}

// SendMsg - Send telegram message
func SendMsg(response tgbotapi.Chattable) {
	if _, err := bot.Send(response); err != nil {
		log.Panicln(err)
	}
}

func breakCommand(message string) (string, []string, bool) {
	var command []string
	var arguments []string
	if message == "" {
		return "", arguments, false
	}

	command = strings.Split(message, " ")
	if len(command) >= 2 {
		arguments = strings.Split(command[1], ",")
	}

	return command[0], arguments, true
}

// StartCommand - Command
func StartCommand(update tgbotapi.Update) {
	stories := []Message{
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "Hi 🙂!", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "My name is Luc", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "... or rather his small digital copy!", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "How can I help you?", 
			MsgType: "Message", 
			Keyboard: &homeReplyKeyboard,
		},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

// TornaCommand - Command
func TornaCommand(update tgbotapi.Update) {
	stories := []Message{
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "What are you interested in?", 
			MsgType: "Message", 
			Keyboard: &homeReplyKeyboard,
		},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

// StoryCommand - Command
func StoryCommand(update tgbotapi.Update) {
	stories := []Message{
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "Ok!", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "This is me:", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "assets/luc.jpg", 
			MsgType: "NewPhotoUpload",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "As I said my name is * Luc Michalski *, I am * 40 years old * and I live in Lyon (France).", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "I love the skydiving, savate (french boxe) but above all * my job *!", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "I currently work at Eedama where I hold the role * Senior Backend Developer * and I work daily with these technologies:", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "In the current workplace I find myself working daily with the following technologies:", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "-Go \n- PHP\n- MySQL\n- Docker\n- Docker-Compose\n- Github\n", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "I started my career at Evolutive Business Group in France where I specialized in e-commerce...", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "Then, I moved to UK, where I had the tremedous luck to work for We Are Social, as senior social media technologist, specialized in social networks (Facebook, Twitter, Instagram)...", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "My second work experience in UK was Blippar, company that allowed my to get a visa L1-A to New York, where I operated as the Global Head of Server", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "This is only a part of my knowledge, for the list and a complete detail you can use the button * Technologies and Projects * in the menu below.", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "In the * Contact * section I leave you the link to my * GitHub * where you can check the quality of my code, such as this Bot, without wasting time doing those boring and useless tests. NoTest * #! *", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 2, 
			Content: "If you think my figure can be useful for your project and if you have an interesting proposal, feel free to contact me!", 
			MsgType: "Message",
		},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

// JobsCommand - Command
func JobsCommand(update tgbotapi.Update) {
	stories := []Message{
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "*Senior PHP Backend & Rest API Developer*\nGiugno 2018 - OGGI\n*Facile.it S.p.A*\n\nAgency; Sviluppo e mantenimento Web Application e servizi Rest API \n\nTecnologie usate/apprese: PHP - MySQL - Symfony - Docker - k8s - GitLab - Redis - Kibana - RabbitMQ", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "*Senior PHP Backend & Rest API Developer*\nAprile 2015 - Giugno 2018\n*S2K Agency*\n\nAgency; Sviluppo e mantenimento Web Application e servizi Rest API \n\nTecnologie usate/apprese: PHP - MySQL - Laravel - Docker - Git - Redis - Deployer", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "*Junior PHP Web Developer*\nMaggio 2014 - Marzo 2015\n*Pro Web Consulting*\n\nAgency; Sviluppo e mantenimento Web Application.\n\nTecnologie usate/apprese: PHP - MySQL - Laravel - Homestead - Git", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "*Junior PHP Web Developer*\nFebbraio 2012 - Aprile 2014\n*Touring Club Italiano*\n\nSviluppo e mantenimento dei canali pubblici principali di Touring Club Italiano e Bandiere Arancioni.\n\nTecnologie usate/apprese: PHP - MySQL - CodeIgniter - Drupal", 
			MsgType: "Message",
		},
		Message{
			ChatID: update.Message.Chat.ID, 
			Duration: 1, 
			Content: "*Tester Funzionale, PMO*\nOttobre 2011 - Febbraio 2012\n*NTT DATA Italia*\n\nMi occupavo principalmente di eseguire dei test funzionali su applicativi riguardanti la pubblicazione e gestione pubblicità a livello web, stampa e radio per il *GRUPPO SOLE 24 ORE*.", 
			MsgType: "Message",
		},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

// TechCommand - Command
func TechCommand(update tgbotapi.Update) {
	stories := []Message{
		// Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Languages*: \n\n -*PHP* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Go* ⭐️️️️⭐️⭐️⭐️ \n -*Python* ⭐️️️️⭐️⭐️\n -*C#* ⭐️️️️⭐️\n -*Rust* ⭐️️️️⭐️ ", MsgType: "Message"},
		// Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Database*: \n\n -*MySQL* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*MongoDB* ⭐️️️️⭐️⭐️ \n", MsgType: "Message"},
		// Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Framework*: \n\n -*Symfony* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Laravel* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Codeigniter* ⭐️️️️⭐️⭐ \n -*Rocket* ⭐️️️️⭐️", MsgType: "Message"},
		// Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Cache*: \n\n -*Redis* ⭐️️️️⭐️⭐️⭐️ ", MsgType: "Message"},
		// Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Other*: \n\n -*Docker* ⭐️️️️⭐️⭐️⭐️ \n -*RabbitMQ* ⭐️️️️⭐️⭐ \n -*k8s* ⭐️️️️", MsgType: "Message"},
		Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Languages*: \n -*PHP* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Go* ⭐️️️️⭐️⭐️⭐️ \n -*Python* ⭐️️️️⭐️⭐️\n -*C#* ⭐️️️️⭐️\n -*Rust* ⭐️️️️⭐️ \n\n*Database*: \n -*MySQL* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*MongoDB* ⭐️️️️⭐️⭐️ \n\n*Framework*: \n -*Symfony* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Laravel* ⭐️️️️⭐️⭐️⭐️⭐️  \n -*Codeigniter* ⭐️️️️⭐️⭐ \n -*Rocket* ⭐️️️️⭐️\n\n*Cache*: \n -*Redis* ⭐️️️️⭐️⭐️⭐️ \n\n*Altro*: \n -*Docker* ⭐️️️️⭐️⭐️⭐️ \n -*RabbitMQ* ⭐️️️️⭐️⭐ \n -*k8s* ⭐️️️️", MsgType: "Message"},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

// ContactsCommand - Command
func ContactsCommand(update tgbotapi.Update) {
	stories := []Message{
		Message{ChatID: update.Message.Chat.ID, Duration: 1, Content: "*Email*: michalski.luc@gmail.com \n*Linkedin*: https://www.linkedin.com/in/luc-m-2751909/ \n*Github*: https://github.com/lucmichalski", MsgType: "Message"},
	}

	for _, story := range stories {
		consumeChainMessage(story)
	}
}

func main() {
	// Bot
	log.Println("Start Bot")
	updatesHandler()
}
