package twitch

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brunobolting/go-twitch-chat/ws"
	"github.com/gorilla/websocket"
)

type Twitch struct {
	conn      *websocket.Conn
	send      chan []byte
	Message   chan Message
	close     chan interface{}
	interrupt chan os.Signal
	Answer chan Message
}

func dial() *websocket.Conn {
	socketUrl := "wss://irc-ws.chat.twitch.tv:443"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	log.Println("connected to IRC server...")

	return conn
}

func NewTwitch() *Twitch {
	tw := &Twitch{
		conn:      dial(),
		send:      make(chan []byte),
		Message:   make(chan Message),
		close:     make(chan interface{}),
		interrupt: make(chan os.Signal),
		Answer:    make(chan Message),
	}

	// signal.Notify(tw.interrupt, os.Interrupt)

	return tw
}

func (tw *Twitch) Run(channel string, hub *ws.Hub) {
	defer tw.conn.Close()

	tw.handleshake(channel)

	go tw.readPump()
	go tw.writePump()
	go tw.handleMessage(hub)

	<-make(chan struct{})
	return
}

func (tw *Twitch) handleshake(channel string) {
	// tw.send <- []byte("CAP REQ twitch.tv/tags")
	err := tw.conn.WriteMessage(websocket.TextMessage, []byte("CAP REQ twitch.tv/tags"))
	if err != nil {
		log.Fatal(err)
	}

	nickname := generateAnonymousNickname()

	// tw.send <- []byte("NICK " + nickname)
	err = tw.conn.WriteMessage(websocket.TextMessage, []byte("NICK "+nickname))
	if err != nil {
		log.Fatal(err)
	}

	// tw.send <- []byte("JOIN #" + channel)
	err = tw.conn.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channel))
	if err != nil {
		log.Fatal(err)
	}
}

func generateAnonymousNickname() string {
	rand.Seed(time.Now().UnixNano())
	randomFiveDigitNumber := ""
	for i := 0; i < 5; i++ {
		randomFiveDigitNumber += (string)(rand.Intn(10) + 48)
	}

	return "justinfan" + randomFiveDigitNumber
}

func (tw *Twitch) readPump() {
	defer close(tw.close)

	for {
		_, msg, err := tw.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		tw.Message <- parse(string(msg))
	}
}

func (tw *Twitch) writePump() {
	for {
		select {
		case message := <-tw.send:
			err := tw.conn.WriteMessage(websocket.TextMessage, message)
			if err == nil {
				log.Printf("send: %s", message)
			}
		case <-tw.interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := tw.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-tw.close:
				log.Println("Channels Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

func (tw *Twitch) handleMessage(hub *ws.Hub) {
	defer close(tw.close)

	for {
		select {
		case message := <-tw.Message:
			if message.Command.Command == "PING" {
				tw.send <- []byte("PONG")
			}

			if message.Command.Command == "PRIVMSG" {
				tw.Answer <- message
				// log.Printf("%s | %s: %s\n", message.Command.Command, message.Author, message.Message)
				// json, _ := json.Marshal(message)
				// log.Printf(string(json))
				// hub.Broadcast <- json
			}
		}
	}
}
