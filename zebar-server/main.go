package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	monitor := NewMonitor(ctx)
	ru := NewResinUpdater()

	go monitor.Run()
	defer monitor.Stop()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, monitor, ru)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Server is shutting down...")
	<-ctx.Done()
	log.Println("Server exited properly")
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request, m *Monitor, u *ResinUpdater) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	go u.RunDailNoteUpdates(conn, GenshinConfig)
	go u.RunDailNoteUpdates(conn, StarRailConfig)
	go u.RunDailNoteUpdates(conn, ZZZConfig)

	listen := m.Register()
	defer m.Unregister(listen)

	for {
		event := <-listen

		if event.Type == StopEvent {
			log.Println("sending after stop event")
			switch event.Name {
			case GenshinProcess:
				go u.RunDailNoteUpdates(conn, GenshinConfig)
			case StarRailProcess:
				go u.RunDailNoteUpdates(conn, StarRailConfig)
			case ZZZProcess:
				go u.RunDailNoteUpdates(conn, ZZZConfig)
			default:
				continue
			}
		}
	}
}
