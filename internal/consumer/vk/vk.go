package vk

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"tgProdLoader/internal/models"

	"github.com/SevereCloud/vksdk/api/params"
	"github.com/SevereCloud/vksdk/v3/api"
)

type vkConsumer struct {
	Client  *api.VK
	GroupID int
	log     *slog.Logger
}

func New(log *slog.Logger, c *api.VK, vkGroupID int) *vkConsumer {
	return &vkConsumer{
		Client:  c,
		GroupID: vkGroupID,
		log:     log,
	}
}

func (vk *vkConsumer) Load(log *slog.Logger, prodChan chan models.Product) {
	vk.log.With("op", "internal.consumer.vk.Load")

	for p := range prodChan {

		vk.log.Debug("product", "name", p.Name)

		mainPhoto, err := vk.loadMainPhoto(p.MainPictureURL)
		if err != nil {
			vk.log.Error("Error:", "can't load mainPhoto:", err.Error())
		}
		vk.log.Debug("loaded mainPhoto to VK", "name", p.Name)

		idUploadedPhotos, err := vk.loadPhotosFromUrls(p.PicturesURL)
		if err != nil {
			vk.log.Error("Error:", "can't load photo:", err.Error())
		}
		vk.log.Debug("loaded Photos to VK", "name", p.Name)

		pars := params.NewMarketAddBuilder()
		vk.log.Debug("loaded mainPhoto to VK", "name", p.Name)

		pars.OwnerID(-vk.GroupID)
		pars.Name(p.Name)
		pars.MainPhotoID(mainPhoto[0].ID)
		pars.PhotoIDs(idUploadedPhotos)
		pars.Description(p.Description)
		pars.Price(float64(p.Price))
		pars.CategoryID(p.CategoryID)

		response, err := vk.Client.MarketAdd(api.Params(pars.Params))
		if err != nil {
			vk.log.Error("can't load Product to VK", "err", err.Error())
		}
		log.Info("Product loaded to VK", "name", p.Name, "productID", response.MarketItemID)

		continue

	}

}

func (vk *vkConsumer) loadPhotosFromUrls(urls []string) ([]int, error) {
	files := make([]io.Reader, 0)

	for _, url := range urls {
		r, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("can't get photo:%s", err.Error())
		}

		files = append(files, io.LimitReader(r.Body, 5242880))

	}

	idPhotos := make([]int, 0)
	for _, pic := range files {
		if len(idPhotos) == 4 {
			break
		}

		resp, err := vk.Client.UploadMarketPhoto(vk.GroupID, false, pic)
		if err != nil {
			log.Printf("Can't upload market photo: %e", err)
		}

		if len(resp) >= 1 {
			idPhotos = append(idPhotos, resp[0].ID)
		}
	}

	return idPhotos, nil
}

func (vk *vkConsumer) loadMainPhoto(url string) (api.PhotosSaveMarketPhotoResponse, error) {
	vk.log.With("+op", "loadMainPhoto")

	req, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't http.GET MainPhoto: %s", err.Error())
	}

	file := io.LimitReader(req.Body, 5242880)

	photoRespVk, err := vk.Client.UploadMarketPhoto(vk.GroupID, true, file)
	if err != nil {
		return nil, fmt.Errorf("can't load MainPhoto to VK: %s", err.Error())
	}

	if len(photoRespVk) == 0 {
		return nil, fmt.Errorf("len of mainPhotos < 1")
	}

	return photoRespVk, nil
}
