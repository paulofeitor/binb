package database

import (
	"fmt"

	"github.com/paulofeitor/binb/internal/pkg/config"
	"github.com/paulofeitor/binb/internal/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client interface {
	GetRooms() ([]*model.Room, error)
}

type client struct {
	conf config.Configuration
	db   *gorm.DB
}

func New(c config.Configuration) (*client, error) {
	db := &client{
		conf: c,
	}

	err := db.connect()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *client) connect() (err error) {
	fmt.Println("DSN:", c.dsn())
	c.db, err = gorm.Open(mysql.Open(c.dsn()), &gorm.Config{})
	if err != nil {
		// TODO :: add log
		fmt.Println("error on database:", err)
		return err
	}
	return nil
}

func (c *client) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.conf.Database.User,
		c.conf.Database.Pass,
		c.conf.Database.Host,
		c.conf.Database.Name,
	)
}

func (c *client) GetRooms() ([]*model.Room, error) {
	var rooms []*model.Room
	if err := c.db.Model(&model.Room{}).
		Where("enabled = ?", false).
		Find(&rooms).
		Error; err != nil {
		// TODO :: add log
		fmt.Printf("database query error: %s\n", err)
		return nil, err
	}
	return rooms, nil
}
