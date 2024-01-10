package ws

import (
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

var msgs = make(chan string, 100)

func init() {
	go func() {
		for {
			msg := <-msgs
			log.Printf("init msg %s\n", msg)

			for _, c := range subscribers {
				select {
				case c.msgs <- msg:
				default:
					c.disconnect()
				}
			}
		}
	}()
}

func Publish(w http.ResponseWriter, r *http.Request) {
	msgs <- r.FormValue("chat_message")

	w.Write([]byte(`<input autofocus name="chat_message" placeholder="say something">`))
}

type connection struct {
	msgs       chan string
	disconnect func()
}

var subscribers = make(map[*http.Request]connection)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	subscribers[r] = connection{
		msgs: make(chan string, 10),
		disconnect: func() {
			fmt.Println("connection disconnecting")
			delete(subscribers, r)
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}

	msgs <- "connected"

	for {
		select {
		case msg := <-subscribers[r].msgs:
			text := fmt.Sprintf(`<div id="notifications" hx-swap-oob="beforeend"><div>%s</div></div>`, msg)
			err := c.Write(ctx, websocket.MessageText, []byte(text))
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			subscribers[r].disconnect()
			msgs <- "disconnected"
			return
		}
	}
}
