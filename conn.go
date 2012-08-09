package main

import (
  "fmt"
  "encoding/json"
	"code.google.com/p/go.net/websocket"
)

type connection struct {
	ws *websocket.Conn
  send chan Event
}

type Event struct {
  Name string `json:"eventName"`
  Data interface{} `json:"data"` // takes arbitrary datas
}

func (c *connection) reader() {
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			break
		}

    var ev Event
    if err := json.Unmarshal([]byte(message), &ev); err != nil {
      fmt.Println(err)
      break
    }
    h.broadcast <- ev
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for ev := range c.send {
    // server should add color to each chat message
    // we want this: {msg: "im saying something", color: "#aasdi"}

    msg, err := json.Marshal(ev)
    if err != nil {
      break
    }

    //check ev.Name, if it's == 'chat msg'... this becomes a real handler.
    // so if it is not, just send it, but if it is
    // err wait. this feels a bit backwards no?
    // like why are we adjusting the server to the client?
    // should be the other way around?
    // if chat msg is arbitrary that is
    // how do u send the color to the client?
    // i guess, but that's... application specific?

    // yeah it's app specific... is WWWWWWWWWWhat i WWWWWWas talking about
    // with the event handlers stuff
    // need a way to register event handlers and parse the data and shit.
    // bleh....................... BLEHHHHHHHHHHHH BLUH
    // blub blug.
    // this would be so ez in js
    // LOL
    // stfu
    // js nerd.  "register event handdlers and parse the data and shit"
    // string ops in go?
    // how done in js btw?
    if string(msg) == 'chat msg' {
      msg :=  string(msg // ...
    if err := websocket.Message.Send(c.ws, string(msg)); err != nil {
			break
		}
	}
	c.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
  fmt.Print("New websocket connection ", ws)
	c := &connection{send: make(chan Event, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}
