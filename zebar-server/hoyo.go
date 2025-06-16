package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ResinUpdater struct {
	mu      sync.Mutex
	notes   map[GameId]DailyNoteCommon
	cancels map[GameId]context.CancelFunc
}

func NewResinUpdater() *ResinUpdater {
	return &ResinUpdater{
		notes:   make(map[GameId]DailyNoteCommon),
		cancels: make(map[GameId]context.CancelFunc),
	}
}

func (u *ResinUpdater) RunDailyNoteUpdates(conn *websocket.Conn, config GameConfig) error {
	note, err := DailyNote(config)
	if err != nil {
		return err
	}

	go u.Run(conn, note)

	return nil
}

func writeNoteToConn(conn *websocket.Conn, note DailyNoteCommon) error {
	return conn.WriteJSON(struct {
		Curr int    `json:"curr"`
		Max  int    `json:"max"`
		Game string `json:"game"`
	}{
		Curr: note.Current,
		Max:  note.Max,
		Game: string(note.Game),
	})
}

func (ru *ResinUpdater) Run(conn *websocket.Conn, note DailyNoteCommon) {

	ru.mu.Lock()

	if cancel, ok := ru.cancels[note.Game]; ok {
		cancel()
	}

	ru.notes[note.Game] = note
	if err := writeNoteToConn(conn, note); err != nil {
		ru.mu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ru.cancels[note.Game] = cancel

	ru.mu.Unlock()

	if note.Current >= note.Max {
		return
	}

	rem := note.FullyRecoveredTs % int(note.RecoverInterval.Seconds())
	time.Sleep(time.Duration(rem) * time.Second)

	ticker := time.NewTicker(note.RecoverInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ru.mu.Lock()

			n, ok := ru.notes[note.Game]
			if !ok {
				ru.mu.Unlock()
				return
			}

			n.Current += 1
			ru.notes[note.Game] = n
			if err := writeNoteToConn(conn, n); err != nil {
				ru.mu.Unlock()
				return
			}

			if n.Current >= n.Max {
				ru.mu.Unlock()
				return
			}

			ru.mu.Unlock()

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func buildRequest(config GameConfig) (req *http.Request, err error) {
	ds := generateDS()

	switch config.game {
	case ZZZ:
		req, err = http.NewRequest("GET", "https://sg-public-api.hoyolab.com/event/game_record_zzz/api/zzz/note?server=prod_gf_us&role_id=1000482805", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("DS", ds)
		req.Header.Set("Cookie", config.cookie)
		req.Header.Set("x-rpc-page", "v1.7.1_#/zzz")
		req.Header.Set("x-rpc-geetest_ext", `{"viewUid":"33046672","server":"prod_gf_us","gameId":8,"page":"v1.7.1_#/zzz","isHost":1,"viewSource":1,"actionSource":127}`)
		req.Header.Set("x-rpc-client_type", "5")
		req.Header.Set("x-rpc-language", "en-us")
		req.Header.Set("User-Agent", "Mozilla/5.0")
	case STARRAIL, GENSHIN:
		url := fmt.Sprintf(
			"https://bbs-api-os.hoyolab.com/game_record/%s/api/%s?role_id=%s&server=%s",
			config.gamePath,
			config.path,
			config.uid,
			config.server,
		)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("DS", ds)
		req.Header.Set("Cookie", config.cookie)
		req.Header.Set("x-rpc-client_type", "5")
		req.Header.Set("x-rpc-language", "en-us")
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("x-rpc-app_version", config.version)

	}
	return req, nil
}

func DailyNote(config GameConfig) (DailyNoteCommon, error) {

	req, err := buildRequest(config)
	if err != nil {
		return DailyNoteCommon{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DailyNoteCommon{}, err
	}
	defer resp.Body.Close()

	log.Println("Fetched note for ", config.game, " status: ", resp.Status)

	note := DailyNoteCommon{
		Game:            config.game,
		RecoverInterval: config.resinRecharge,
	}

	switch config.game {
	case GENSHIN:
		var result DailyNoteResponseGenshin
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		note.Current = result.Data.CurrentResin
		note.Max = result.Data.MaxResin
	case STARRAIL:
		var result DailyNoteResponseStarRail
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		note.Current = result.Data.CurrentStamina
		note.Max = result.Data.MaxStamina
	case ZZZ:
		var result DailyNoteResponseZZZ
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return DailyNoteCommon{}, err
		}
		note.Current = result.Data.Energy.Progress.Current
		note.Max = result.Data.Energy.Progress.Max
	}

	note.RecoverInterval = config.resinRecharge

	return note, nil
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
