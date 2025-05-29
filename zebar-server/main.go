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

	go monitor.Run()
	defer monitor.Stop()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, monitor)
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

func serveWs(w http.ResponseWriter, r *http.Request, m *Monitor) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

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
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
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
