package models

type Product struct {
	Name           string
	Description    string
	Size           string
	Status         string
	Price          int
	MediaGroupID   string
	MainPictureURL string
	PicturesURL    []string
	CategoryID     int
}
