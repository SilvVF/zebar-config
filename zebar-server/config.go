package main

import (
	"flag"
)

type GameId = string

var (
	addr = flag.String("addr", "localhost:45456", "http service address")
)

const (
	GENSHIN        GameId = "genshin"
	STARRAIL       GameId = "hkrpg"
	ZZZ            GameId = "zzz"
	uidGenshin            = "604392290"
	uidHSR                = "614963011"
	serverGenshin         = "os_usa"
	serverHsr             = "prod_official_usa"
	versionGenshin        = "2.11.1"
	versionHsr            = "2.50.1"
	dsSalt                = "6s25p5ox5y14umn1p61aqyyvbvvl3lrt"
)

type GameConfig struct {
	game     GameId
	path     string
	uid      string
	gamePath string
	server   string
	cookie   string
	version  string
}

var ZZZConfig = GameConfig{
	game:     ZZZ,
	gamePath: "zzz",
	path:     "note",
	uid:      "1000482805",
	server:   "prod_gf_us",
	cookie:   HOYOLAB_COOKIE_STARRAIL,
	version:  "",
}

var StarRailConfig = GameConfig{
	game:     STARRAIL,
	gamePath: "hkrpg",
	path:     "note",
	uid:      uidHSR,
	server:   serverHsr,
	cookie:   HOYOLAB_COOKIE_STARRAIL,
	version:  versionHsr,
}

var GenshinConfig = GameConfig{
	game:     GENSHIN,
	gamePath: "genshin",
	path:     "dailyNote",
	uid:      uidGenshin,
	server:   serverGenshin,
	cookie:   HOYOLAB_COOKIE_GENSHIN,
	version:  versionGenshin,
}
