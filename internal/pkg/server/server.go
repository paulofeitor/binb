package server

import (
	"fmt"
	"net/http"

	"github.com/paulofeitor/binb/internal/pkg/config"
	"github.com/paulofeitor/binb/internal/pkg/database"
	"github.com/paulofeitor/binb/internal/pkg/websocket"
)

type Server interface {
	Start() error
}

type server struct {
	c    config.Configuration
	db   database.Client
	pool *websocket.Pool
}

func New(c config.Configuration, db database.Client) *server {
	return &server{
		c:    c,
		db:   db,
		pool: websocket.NewPool(),
	}
}

func (s *server) setupRoutes() {
	http.HandleFunc("/ws", s.serveWs)
}

func (s *server) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
		Conn: ws,
		Pool: s.pool,
	}

	s.pool.Register <- client
	client.Read()
}

func (s *server) Start() {
	s.setupRoutes()
	go s.pool.Start()

	fmt.Println("Distributed Chat App v0.01")
	http.ListenAndServe(":8080", nil)
}
