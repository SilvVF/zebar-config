package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mu sync.Mutex
var listeners = map[chan MonitorEvent]struct{}{}

func register() chan MonitorEvent {
	ch := make(chan MonitorEvent, 1)
	mu.Lock()
	listeners[ch] = struct{}{}
	mu.Unlock()
	return ch
}

func unregister(ch chan MonitorEvent) {
	mu.Lock()
	delete(listeners, ch)
	mu.Unlock()
	close(ch)
}

func main() {
	flag.Parse()

	ctx := context.Background()

	monitor := NewMonitor(ctx)

	go monitor.Run()
	defer monitor.Stop()

	go func() {
		for {
			event := <-monitor.events
			fmt.Printf("Received MonitorEvent %v \n", event)
			for l := range listeners {
				l <- event
			}
		}
	}()

	http.HandleFunc("/ws", serveWs)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	sendDailyNote := func(config GameConfig) {
		resin, max, err := DailyNote(config)
		fmt.Println("sending ", config.game, resin, max, "err: ", err)
		if err != nil {
			return
		}

		if err := conn.WriteJSON(struct {
			Curr int    `json:"curr"`
			Max  int    `json:"max"`
			Game string `json:"game"`
		}{
			Curr: resin,
			Max:  max,
			Game: string(config.game),
		}); err != nil {
			return
		}
	}

	sendDailyNote(GenshinConfig)
	sendDailyNote(StarRailConfig)
	sendDailyNote(ZZZConfig)

	listen := register()
	defer unregister(listen)

	for {
		event := <-listen

		if event.Type == StopEvent {
			fmt.Println("sending after stop event")
			switch event.Name {
			case GenshinProcess:
				sendDailyNote(GenshinConfig)
			case StarRailProcess:
				sendDailyNote(StarRailConfig)
			case ZZZProcess:
				sendDailyNote(ZZZConfig)
			default:
				continue
			}
		}
	}
}

func dailyNoteZZZ() (int, int, error) {
	config := ZZZConfig

	req, err := http.NewRequest("GET", "https://sg-public-api.hoyolab.com/event/game_record_zzz/api/zzz/note?server=prod_gf_us&role_id=1000482805", nil)
	if err != nil {
		return -1, -1, err
	}

	req.Header.Set("DS", generateDS())
	req.Header.Set("Cookie", config.cookie)
	req.Header.Set("x-rpc-page", "v1.7.1_#/zzz")
	req.Header.Set("x-rpc-geetest_ext", `{"viewUid":"33046672","server":"prod_gf_us","gameId":8,"page":"v1.7.1_#/zzz","isHost":1,"viewSource":1,"actionSource":127}`)
	req.Header.Set("x-rpc-client_type", "5")
	req.Header.Set("x-rpc-language", "en-us")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, err
	}

	var res DailyNoteResponseZZZ
	err = json.Unmarshal(b, &res)
	if err != nil {
		return -1, -1, err
	}

	return res.Data.Energy.Progress.Current, res.Data.Energy.Progress.Max, nil
}

func DailyNote(config GameConfig) (int, int, error) {

	if config.game == ZZZ {
		return dailyNoteZZZ()
	}

	url := fmt.Sprintf(
		"https://bbs-api-os.hoyolab.com/game_record/%s/api/%s?role_id=%s&server=%s",
		config.gamePath,
		config.path,
		config.uid,
		config.server,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, err
	}

	ds := generateDS()

	req.Header.Set("DS", ds)
	req.Header.Set("Cookie", config.cookie)
	req.Header.Set("x-rpc-client_type", "5")
	req.Header.Set("x-rpc-language", "en-us")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("x-rpc-app_version", config.version)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	resin := -1
	max := -1

	if config.game == GENSHIN {
		var result DailyNoteResponseGenshin
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return -1, -1, err
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return -1, -1, err
		}
		resin = result.Data.CurrentResin
		max = result.Data.MaxResin
	} else {
		var result DailyNoteResponseStarRail
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return -1, -1, err
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return -1, -1, err
		}
		resin = result.Data.CurrentStamina
		max = result.Data.MaxStamina
	}
	return resin, max, nil
}

func generateDS() string {
	t := time.Now().Unix()

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	random := make([]byte, 6)
	for i := range 6 {
		random[i] = letters[rand.Intn(len(letters))]
	}

	// Format string to hash
	raw := fmt.Sprintf("salt=%s&t=%d&r=%s", dsSalt, t, string(random))
	hash := fmt.Sprintf("%x", md5.Sum([]byte(raw)))

	return fmt.Sprintf("%d,%s,%s", t, string(random), hash)
}
