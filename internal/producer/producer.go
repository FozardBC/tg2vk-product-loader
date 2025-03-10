package producer

import (
	"log/slog"
	"tgProdLoader/internal/models"
)

type Producer interface {
	HandleMessages(log *slog.Logger, product chan models.Product) error
}
