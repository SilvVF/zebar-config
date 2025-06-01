package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	ytmusic *websocket.Conn

	conns map[*websocket.Conn]struct{}
	last  []byte
	mu    sync.Mutex
}

func (s *Server) TogglePlayback() {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println("checking yt conn")
	if s.ytmusic != nil {
		log.Println("sending toggle playback to yt-music")
		s.ytmusic.WriteMessage(websocket.TextMessage, make([]byte, 0))
	}
}

func (s *Server) Remove(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.conns, conn)
}

func (s *Server) Add(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.last) != 0 {
		conn.WriteMessage(websocket.TextMessage, s.last)
	}
	s.conns[conn] = struct{}{}
}

func (s *Server) Broadcast(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.last = data

	for conn := range s.conns {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func NewServer() *Server {
	return &Server{
		last:  make([]byte, 0),
		conns: map[*websocket.Conn]struct{}{},
	}
}

func main() {
	flag.Parse()

	ctx := context.Background()

	monitor := NewMonitor(ctx)

	go monitor.Run()
	defer monitor.Stop()

	serv := NewServer()

	http.HandleFunc("/ytmusic", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		serv.ytmusic = conn

		for {
			mtype, b, err := conn.ReadMessage()
			if err != nil {
				return
			}
			switch mtype {
			case websocket.TextMessage:
				log.Println("recieved ytmusic msg: ", string(b))
				serv.Broadcast(b)
			}
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, monitor, serv)
	})

	serverError := make(chan error, 1)

	go func() {
		log.Printf("Server is running on http://localhost%s", *addr)
		if err := http.ListenAndServe(*addr, nil); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Println("Server is shutting down...")
	<-ctx.Done()
	log.Println("Server exited properly")
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request, m *Monitor, s *Server) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.Add(conn)

	defer conn.Close()
	defer s.Remove(conn)

	u := NewResinUpdater()

	go u.RunDailyNoteUpdates(conn, GenshinConfig)
	go u.RunDailyNoteUpdates(conn, StarRailConfig)
	go u.RunDailyNoteUpdates(conn, ZZZConfig)

	listen := m.Register()
	defer m.Unregister(listen)

	done := make(chan struct{}, 1)

	go func() {

		defer close(done)

		for {
			mtype, b, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				break
			}

			switch mtype {
			case websocket.TextMessage:
				log.Println("recieved ", string(b))
				if string(b) == "toggle-playback" {
					s.TogglePlayback()
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case event, ok := <-listen:
			if !ok {
				return
			}

			if event.Type == StopEvent {
				log.Println("sending after stop event")
				switch event.Name {
				case GenshinProcess:
					go u.RunDailyNoteUpdates(conn, GenshinConfig)
				case StarRailProcess:
					go u.RunDailyNoteUpdates(conn, StarRailConfig)
				case ZZZProcess:
					go u.RunDailyNoteUpdates(conn, ZZZConfig)
				default:
					continue
				}
			}
		}
	}
}
