package models

import "github.com/jinzhu/gorm"

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

type ServicesConfig func(*Services) error

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

//func NewServices(dialect, connStr string) (*Services, error) {
//	db, err := gorm.Open(dialect, connStr)
//	if err != nil {
//		return nil, err
//	}
//	db.LogMode(true)
//	return &Services{
//		User:    NewUserService(db),
//		Gallery: NewGalleryService(db),
//		Image:   NewImageService(),
//		db:      db,
//	}, nil
//}

// Close method closes the database connection.
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate method will attempt to automatically migrate all tables.
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// DestructiveReset method drops all tables and rebuilds them.
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
