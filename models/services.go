package models

import "github.com/jinzhu/gorm"

func NewServices(connStr string) (*Services, error) {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: &galleryGorm{},
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
