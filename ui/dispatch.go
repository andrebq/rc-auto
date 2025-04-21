package ui

import (
	"encoding/json"
	"net/http"
	"rcauto/bus"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type ()

var connections uint64

func Dispatch(pump bus.Pump) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		id := atomic.AddUint64(&connections, 1)
		defer func() {
			c.Close()
		}()
		pump.Open(id)
		var wg sync.WaitGroup
		wg.Add(2)
		handleInput := func() {
			defer func() {
				wg.Done()
				pump.Close(id)
			}()
			defer c.Close()
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					return
				}
				if mt == websocket.TextMessage {
					var pl bus.Payload
					err := json.Unmarshal(message, &pl)
					if err == nil {
						pump.Push(id, pl)
					}
				}
			}
		}
		handleOutput := func() {
			defer wg.Done()
			for {
				msgs, err := pump.Pop(id)
				if err != nil {
					return
				}
				for _, m := range msgs {
					message, err := json.Marshal(m)
					if err != nil {
						continue
					}
					err = c.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						break
					}
				}
			}
		}
		go handleInput()
		go handleOutput()
		wg.Wait()
	})
}
