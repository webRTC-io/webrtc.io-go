package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"github.com/oskarth/wsevents"
	"log"
	"net/http"
	"text/template"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("index.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func resourceHandler(c http.ResponseWriter, req *http.Request) {
	fmt.Print("Serving "+req.URL.Path[1:], "\n")
	http.ServeFile(c, req, req.URL.Path[1:])
}

func main() {
	flag.Parse()
	wsevents.Connect(func(c *wsevents.Connection) {

		// wrap event with color before sending it
		c.On("chat msg", func(msg interface{}) {
			d := map[string]string{
				"msg":   msg.(string),
				"color": "#999",
			}
			c.Broadcast("chat msg", d)
		})
	})
	go wsevents.Run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/resources/webrtc.io.js", resourceHandler)
	http.HandleFunc("/resources/style.css", resourceHandler)
	http.Handle("/ws", websocket.Handler(wsevents.Handler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
