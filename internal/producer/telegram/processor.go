package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"tgProdLoader/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *tgProducer) loadProducts(prodChan chan models.Product, updates tgbotapi.UpdatesChannel) {
	product := models.Product{}
	prods := make([]models.Product, 0)
OUT:
	for u := range updates {
		if u.Message == nil {
			continue
		}

		if u.Message.MediaGroupID != product.MediaGroupID {
			prods = append(prods, product)
			product.PicturesURL = make([]string, 0)
			product.MediaGroupID = u.Message.MediaGroupID
		}

		switch u.Message.Text {
		case reload:
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
			setKeyboard(&msg, startKeyboard, "Главное меню:")
			t.bot.Send(msg)
			return

		case stopLoad:
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
			setKeyboard(&msg, startKeyboard, "Загрузка прекращена.")
			t.bot.Send(msg)
			t.log.Debug("User stoped loading", "User", u.Message.From.UserName)
			return
		case category:
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")

			t.changeCategory(*u.Message, updates)

			setKeyboard(&msg, activeLoadingKeyboard, "Вы сменили категорию: "+t.category+"Можете продолжить загрузку товаров")
			t.bot.Send(msg)

			continue
		case goLoadServices: // ВЫХОД ИЗ ВНЕШНЕГО ЦИКЛА
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Началась загрузка товаров в ВК. Ожидайте...")
			t.bot.Send(msg)
			t.log.Debug("started to write in channel")
			break OUT
		}

		if u.Message.Text == "" {
			if !u.Message.ForwardFromChat.IsChannel() {
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
				setKeyboard(&msg, startLoading, "Сообщение не является товаром. ОНО ДОЛЖНО БЫТЬ ПЕРЕСЛАНО ИЗ КАНАЛА")
				t.bot.Send(msg)
				return
			}
		}

		if u.Message.Caption != "" {
			product.MediaGroupID = u.Message.MediaGroupID
			getProdInfo(&product, u.Message.Caption)

			t.log.Debug("Added data in struct", "name Prod", product.Name)

			MainPictureURL, err := t.bot.GetFileDirectURL(u.Message.Photo[3].FileID)
			if err != nil {
				log.Print(err)
			}
			t.log.Debug("Added mainPicture in struct", "name Prod", product.Name)

			product.MainPictureURL = MainPictureURL
		} else if product.MediaGroupID == u.Message.MediaGroupID && u.Message.MediaGroupID != "" {
			photo, err := t.bot.GetFileDirectURL(u.Message.Photo[3].FileID)
			if err != nil {
				t.log.Error("can't get direct URL", "err", err.Error())
			}
			product.PicturesURL = append(product.PicturesURL, photo)
			t.log.Debug("Added pictures in struct", "name Prod", product.Name)
		}

	}

	for i, p := range prods {

		if i == 0 { //первый продукт записывается пустым ХЗ особенность логики
			continue
		}

		t.log.Debug("Written in channel", "name Prod", p.Name)
		prodChan <- p
		continue
	}

	seconds := len(prods) * 3 // примерное время ожидания загрузки товара из расчет 1 товар - 3 секунды

	duration := time.Duration(seconds) * time.Second

	time.Sleep(duration)
}

func getProdInfo(p *models.Product, text string) {
	texts := strings.Split(text, "\n")

	p.Name = texts[0]
	p.Size = texts[1]
	p.Status = texts[2]

	for _, s := range texts {
		if strings.HasPrefix(s, "Цена") {
			_, priceStr, _ := strings.Cut(s, " ")
			priceStr = strings.TrimSpace(priceStr)
			priceStr = strings.ReplaceAll(priceStr, ".", "")
			price, err := strconv.Atoi(priceStr)
			if err != nil {
				fmt.Print(err)
			}
			p.Price = price
			break
		}
	}

	p.Description = truncateText(text)

}

func truncateText(text string) string { // Функция для вк чтобы после "свзяь" все обрезать
	index := strings.Index(text, "Связь")
	if index != -1 {
		return text[:index]
	}
	return text
}
