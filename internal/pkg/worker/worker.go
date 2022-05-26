package worker

import (
	"fmt"

	"github.com/paulofeitor/binb/internal/pkg/config"
	"github.com/paulofeitor/binb/internal/pkg/database"
)

type worker struct {
	conf config.Configuration
	db   database.Client
}

func New(c config.Configuration, db database.Client) (*worker, error) {
	return &worker{
		conf: c,
		db:   db,
	}, nil
}

func (w *worker) Start() error {
	rooms, err := w.db.GetRooms()
	if err != nil {
		return err
	}
	for _, r := range rooms {
		fmt.Println(r)
	}
	return nil
}
