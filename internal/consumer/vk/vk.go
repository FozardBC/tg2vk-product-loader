package vk

import (
	"log/slog"
	"sync"
	"tgProdLoader/internal/models"

	"github.com/SevereCloud/vksdk/v3/api"
)

type vkConsumer struct {
	Client *api.VK
}

func New(c *api.VK) *vkConsumer {
	return &vkConsumer{
		Client: c,
	}
}

func (v *vkConsumer) Load(wg *sync.WaitGroup, log *slog.Logger, prodChan chan models.Product) {
	for p := range prodChan {
		log.Debug("product", "name", p.Name)
	}
	wg.Done()
}
