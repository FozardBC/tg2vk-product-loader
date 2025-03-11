package telegram

import (
	"log/slog"
	"tgProdLoader/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgProducer struct {
	bot      *tgbotapi.BotAPI
	category string
	log      *slog.Logger
}

func New(log *slog.Logger, bot *tgbotapi.BotAPI) *tgProducer {
	return &tgProducer{
		bot: bot,
		log: log,
	}
}

func (t *tgProducer) HandleMessages(log *slog.Logger, productChan chan models.Product) error {
	log.Debug("started handling messages")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	log.Info("started telegram bot")

	for u := range updates {

		log.Debug("Handling user:", "User", u.Message.From.UserName) // сделать имя вдеюбаге

		if u.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(u.FromChat().ID, "")

		switch u.Message.Text {
		case reload:
			setKeyboard(&msg, startKeyboard, "Главное меню:")
			t.bot.Send(msg)
			continue

		case "/start":
			log.Debug("User used cmd /start", "User", u.Message.From.UserName)
			setKeyboard(&msg, startKeyboard, "Что хотите сделать?")

			_, err := t.bot.Send(msg)
			if err != nil {
				log.Info("can't send message with startKeyboard:", "error", err.Error())
			}
			continue

		case category:
			t.changeCategory(*u.Message, updates)

			if t.category == "" || t.category == " " {
				t.changeCategory(*u.Message, updates)
			}
			fallthrough
		case startLoad:

			log.Debug("User used cmd startLoad", "User", u.Message.From.UserName)

			if t.category == "" || t.category == " " {
				msg.Text = "Категория не выбрана!"

				t.changeCategory(*u.Message, updates)

				setKeyboard(&msg, startLoading, "Вы выбрали категорию: "+t.category+". Начать загрузку?")

				_, err := t.bot.Send(msg)
				if err != nil {
					log.Info("can't send message with message:", "error", err.Error())
				}

				continue
			}

			//Началась приёмка товаров из сообщений
			setKeyboard(&msg, activeLoadingKeyboard, "Выбранная категория:"+t.category+"\nЗагрзука товаров началась. Пришлите посты.\nПосле того как пришлете все товары выбранной категории - нажмите "+goLoadServices)

			_, err := t.bot.Send(msg)
			if err != nil {
				log.Info("can't send message with startKeyboard:", "error", err.Error())
			}

			log.Debug("Starting loading", "User", u.Message.From.UserName)
			t.loadProducts(productChan, updates)

			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
			setKeyboard(&msg, startKeyboard, "Загрузка товаров завершена. ")
			t.bot.Send(msg)
		default:
			setKeyboard(&msg, startKeyboard, "Неизвестная команда")

			_, err := t.bot.Send(msg)
			if err != nil {
				log.Info("can't send message with startKeyboard:", "error", err.Error())
			}
			continue
		}

	}

	return nil
}

func (t *tgProducer) setCategory(updates tgbotapi.UpdatesChannel) bool {

	t.log.Debug("User used cmd setCategory")

	for u := range updates {
		if u.Message == nil {
			continue
		}

		if u.Message.Text == reload {
			msg := tgbotapi.NewMessage(u.Message.From.ID, "Перезагрузка...")
			t.bot.Send(msg)
			setKeyboard(&msg, startKeyboard, "Главное меню:")
			t.bot.Send(msg)

			return false

		}

		if isCategory(u.Message.Text) {
			t.category = u.Message.Text
			t.log.Debug("User setted category", "User", u.Message.From.UserName, "category", t.category)
			break
		} else {
			msg := tgbotapi.NewMessage(u.Message.From.ID, "Категория не из списка! Выберите категорию из списка!")
			t.bot.Send(msg)
			continue
		}

	}
	return true
}

func (t *tgProducer) changeCategory(m tgbotapi.Message, u tgbotapi.UpdatesChannel) {
	t.log.Debug("User used cmd category", "User", m.From.UserName)

	msg := tgbotapi.NewMessage(m.Chat.ID, "")

	setKeyboard(&msg, chooseKeyboard, "Выберите категорию:")

	_, err := t.bot.Send(msg)
	if err != nil {
		t.log.Info("can't send message with startKeyboard:", "error", err.Error())
	}

	if !t.setCategory(u) {
		t.log.Debug("category isn't chosen. Try agian")
		t.setCategory(u)
	}
}

func isCategory(msg string) bool {
	valid := false

	for _, c := range categories {
		if msg == c {
			valid = true
		}
	}

	return valid
}

func setKeyboard(msg *tgbotapi.MessageConfig, keyboard tgbotapi.ReplyKeyboardMarkup, text string) {

	msg.ReplyMarkup = keyboard
	msg.Text = text

}
