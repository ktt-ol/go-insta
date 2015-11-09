package insta

import (
	"encoding/json"
	"fmt"
	"log"

	_ "net/http/pprof"

	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/tron"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"net/http"
)

var upgrader = websocket.Upgrader{}

type controlMsg struct {
	Player string
	Button string
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error: ws upgrade:", err)
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := controlMsg{}
		if err := json.Unmarshal(message, &msg); err == nil {
			fmt.Println("button:", msg)
		}
		if string(message) == "glider" {
			if s.life != nil {
				log.Println("add glider")
				s.life.AddRandomSpaceship()
			}
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
}

type Server struct {
	life *life.Life
	tron *tron.Game
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetLife(l *life.Life) {
	s.life = l
}

func (s *Server) SetTron(t *tron.Game) {
	s.tron = t
}

func (s *Server) Run() {

	go http.ListenAndServe("localhost:6060", nil)
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/ws", s.serveWs)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":9090", r))
}
