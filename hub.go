package main

import "fmt"

type hub struct {
	connections map[*connection]bool
	broadcast chan Event
	register chan *connection
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan Event),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
      fmt.Print("Register: ", h.connections)
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case ev := <-h.broadcast:
      fmt.Println("Broadcast: ", ev.Name, " : ", ev.Data)
			for c := range h.connections {
				select {
        case c.send <- ev:
				default:
          delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}
