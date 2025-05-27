package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

type GameId = string

var (
	addr = flag.String("addr", "localhost:45456", "http service address")
)

const (
	GENSHIN  GameId = "genshin"
	STARRAIL GameId = "hkrpg"
	ZZZ      GameId = "zzz"

	uidGenshin = "604392290"
	uidHSR     = "614963011"

	serverGenshin = "os_usa"
	serverHsr     = "prod_official_usa"

	versionGenshin = "2.11.1"
	versionHsr     = "2.50.1"

	dsSalt = "6s25p5ox5y14umn1p61aqyyvbvvl3lrt"
)

type GameConfig struct {
	game          GameId
	path          string
	uid           string
	gamePath      string
	server        string
	cookie        string
	version       string
	resinRecharge time.Duration
}

var ZZZConfig = GameConfig{
	game:          ZZZ,
	gamePath:      "zzz",
	path:          "note",
	uid:           "1000482805",
	server:        "prod_gf_us",
	resinRecharge: time.Minute * 6,
}

var StarRailConfig = GameConfig{
	game:          STARRAIL,
	gamePath:      "hkrpg",
	path:          "note",
	uid:           uidHSR,
	server:        serverHsr,
	version:       versionHsr,
	resinRecharge: time.Second * 360,
}

var GenshinConfig = GameConfig{
	game:          GENSHIN,
	gamePath:      "genshin",
	path:          "dailyNote",
	uid:           uidGenshin,
	server:        serverGenshin,
	version:       versionGenshin,
	resinRecharge: time.Second * 480,
}

func init() {
	file, err := os.Open("conf.env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	env := map[string]string{}

	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)

		split := strings.SplitN(line, "=", 2)
		if len(split) != 2 {
			log.Fatal("couldnt parse " + line)
		}
		env[split[0]] = split[1]
	}

	cookie := env["HOYOLAB_COOKIE"]
	ZZZConfig.cookie = cookie
	GenshinConfig.cookie = cookie
	StarRailConfig.cookie = cookie
}
