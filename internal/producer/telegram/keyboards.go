package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	startLoad      = "Начать"
	setings        = "Настройки"
	category       = "Выбрать категорию"
	load           = "Начать загрузку"
	stopLoad       = "Прекратить загрузку"
	goLoadServices = "Начать загрузку в группы"
	reload         = "/reload"
)

var categories = []string{
	"👕 Футболки и поло", "👟 Кроссовки и кеды",
	"👖 Брюки", "🧥 Лёгкие куртки и ветровки",
	"Толстовки и свитшоты", "👜 Сумки",
	"➰ Ремни", "👛 Кошельки", "🧢 Бейсболки",
	"🧢 Шапки", "🧤 Перчатки и варежки",
}

const (
	tshirt    = "👕 Футболки и поло"
	sneakers  = "👟 Кроссовки и кеды"
	pants     = "👖 Брюки"
	jackets   = "🧥 Лёгкие куртки и ветровки"
	sweaters  = "Толстовки и свитшоты"
	bags      = "👜 Сумки"
	belts     = "➰ Ремни"
	portmones = "👛 Кошельки"
	baseballs = "🧢 Бейсболки"
	hats      = "🧢 Шапки"
	gloves    = "🧤 Перчатки и варежки"
)

var activeLoadingKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(goLoadServices),
		tgbotapi.NewKeyboardButton(stopLoad),
		tgbotapi.NewKeyboardButton(category),
	),
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(startLoad),
		tgbotapi.NewKeyboardButton(setings),
	),
)

var startLoading = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(startLoad),
		tgbotapi.NewKeyboardButton(category),
	),
)

var chooseKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(tshirt),
		tgbotapi.NewKeyboardButton(sneakers),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(pants),
		tgbotapi.NewKeyboardButton(jackets),
		tgbotapi.NewKeyboardButton(sweaters),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(bags),
		tgbotapi.NewKeyboardButton(belts),
		tgbotapi.NewKeyboardButton(portmones),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(baseballs),
		tgbotapi.NewKeyboardButton(hats),
		tgbotapi.NewKeyboardButton(gloves),
	),
)
