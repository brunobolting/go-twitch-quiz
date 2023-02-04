package twitch

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Twitch struct {
	conn      *websocket.Conn
	send      chan []byte
	Message   chan Message
	Answer chan Message
	Close chan struct{}
	closeHandler chan struct{}
	wg sync.WaitGroup
}

func dial() *websocket.Conn {
	socketUrl := "wss://irc-ws.chat.twitch.tv:443"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}

	return conn
}

func NewTwitch() *Twitch {
	tw := &Twitch{
		conn:      dial(),
		send:      make(chan []byte),
		Message:   make(chan Message),
		Answer:    make(chan Message),
		Close: make(chan struct{}),
		closeHandler: make(chan struct{}),
	}

	return tw
}

func (tw *Twitch) Run(channel string) {
	defer func() {
		tw.conn.Close()
	}()

	tw.wg.Add(3)

	err := tw.handleshake(channel)
	if err != nil {
		log.Fatalln(err)
	}

	go tw.readPump()
	go tw.writePump()
	go tw.handleMessage()

	tw.wg.Wait()
}

func (tw *Twitch) handleshake(channel string) error {
	err := tw.conn.WriteMessage(websocket.TextMessage, []byte("CAP REQ twitch.tv/tags"))
	if err != nil {
		return err
	}

	nickname := generateAnonymousNickname()

	err = tw.conn.WriteMessage(websocket.TextMessage, []byte("NICK "+nickname))
	if err != nil {
		return err
	}

	err = tw.conn.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channel))
	if err != nil {
		return err
	}

	return nil
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
	defer func()  {
		tw.wg.Done()
		tw.conn.Close()
	}()

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
	defer func()  {
		tw.wg.Done()
		tw.conn.Close()
		close(tw.closeHandler)
	}()

	for {
		select {
		case message := <-tw.send:
			err := tw.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error: %s | send: %s", err, message)
				return
			}
		case <-tw.Close:
			err := tw.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			return
		}
	}
}

func (tw *Twitch) handleMessage() {
	defer func()  {
		tw.wg.Done()
		close(tw.Message)
	}()

	for {
		select {
		case message := <-tw.Message:
			if message.Command.Command == "PING" {
				tw.send <- []byte("PONG")
			}

			if message.Command.Command == "PRIVMSG" {
				tw.Answer <- message
			}
		case <-tw.closeHandler:
			return
		}
	}
}
