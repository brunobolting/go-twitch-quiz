package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/brunobolting/go-twitch-chat/game"
	"github.com/brunobolting/go-twitch-chat/twitch"
	"github.com/brunobolting/go-twitch-chat/usecase/question"
	"github.com/gorilla/websocket"
)

type Message struct {
	Command string `json:"command"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn *websocket.Conn

	game *game.Game

	tw *twitch.Twitch

	interrupt chan struct{}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		close(c.interrupt)
		close(c.tw.Close)
		close(c.game.Close)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		command := &game.Command{}
		json.Unmarshal(msg, &command)
		c.game.Command <- command
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.game.SendToClient:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.game.SendToClient)
			for i := 0; i < n; i++ {
				w.Write(<-c.game.SendToClient)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.interrupt:
			return
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request, service *question.Service) {
	queryParams := r.URL.Query()
	chat := queryParams.Get("channel")
	if chat == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tw := twitch.NewTwitch()

	game := game.NewGame(tw, service)

	client := &Client{conn: conn, game: game, tw: tw, interrupt: make(chan struct{})}

	go tw.Run(chat)
	go game.Run()
	go client.writePump()
	go client.readPump()
}
