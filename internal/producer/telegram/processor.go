package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"tgProdLoader/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *tgProducer) loadProducts(prodChan chan models.Product, updates tgbotapi.UpdatesChannel) {
	product := models.Product{}

	for u := range updates {
		if u.Message == nil {
			continue
		}

		switch u.Message.Text {
		case reload:
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
			setKeyboard(&msg, startKeyboard, "Главное меню:")
			t.bot.Send(msg)
			continue

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

			MainPictureURL, err := t.bot.GetFileDirectURL(u.Message.Photo[3].FileID)
			if err != nil {
				log.Print(err)
			}

			product.MainPictureURL = MainPictureURL
		} else if product.MediaGroupID == u.Message.MediaGroupID && u.Message.MediaGroupID != "" {
			photo, err := t.bot.GetFileDirectURL(u.Message.Photo[3].FileID)
			if err != nil {
				log.Print(err)
			}
			product.PicturesURL = append(product.PicturesURL, photo)
		} else {
			prodChan <- product
			continue
		}

	}

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
