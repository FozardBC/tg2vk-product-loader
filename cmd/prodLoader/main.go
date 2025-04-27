package main

import (
	"log"
	"log/slog"
	"os"
	"sync"
	"tgProdLoader/internal/config"
	"tgProdLoader/internal/consumer/vk"
	"tgProdLoader/internal/lib/logger/handler/slogpretty"
	"tgProdLoader/internal/models"
	"tgProdLoader/internal/producer/telegram"

	"github.com/SevereCloud/vksdk/api/params"
	"github.com/SevereCloud/vksdk/v3/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	tgBot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Error("can't start telegram bot:", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Autharizated telegram bot:", "bot", tgBot.Self.FirstName)

	vkClient := api.NewVK(cfg.VkToken)

	logger.Info("Autharizated vk:", "Name:", getClientName(vkClient))

	var chanProds chan models.Product = make(chan models.Product)

	tgProducer := telegram.New(logger, tgBot)
	vkConsumer := vk.New(logger, vkClient, cfg.VkGroupId)

	wg := sync.WaitGroup{}
	wg.Add(10)

	go tgProducer.HandleMessages(logger, chanProds)
	go vkConsumer.Load(logger, chanProds)

	wg.Wait()

}

func getClientName(vk *api.VK) string {

	p := params.NewAccountGetInfoBuilder()

	info, err := vk.AccountGetProfileInfo(api.Params(p.Params))
	if err != nil {
		log.Fatal(err)
	}

	return info.FirstName + " " + info.LastName
}

// fix it
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
