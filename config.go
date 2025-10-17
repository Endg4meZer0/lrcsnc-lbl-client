package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Template     string             `toml:"template"`
	Protocol     string             `toml:"protocol"`
	ListenAt     string             `toml:"listen-at"`
	InStopAt     int                `toml:"in-stop-at"`
	OutStartAt   int                `toml:"out-start-at"`
	Multiplier   MultiplierConfig   `toml:"multiplier"`
	Instrumental InstrumentalConfig `toml:"instrumental"`
	Errors       ErrorsConfig       `toml:"errors"`
}

type MultiplierConfig struct {
	Enabled bool   `toml:"enabled"`
	Format  string `toml:"format"`
	AddTo   string `toml:"add-to"`
	ToLeft  bool   `toml:"to-left"`
}

type InstrumentalConfig struct {
	Enabled    bool    `toml:"enabled"`
	Interval   float64 `toml:"interval"`
	Symbol     string  `toml:"symbol"`
	MaxSymbols int     `toml:"max-symbols"`
}

type ErrorsConfig struct {
	NotPlaying     string `toml:"not-playing"`
	NoLyrics       string `toml:"no-lyrics"`
	NoSyncedLyrics string `toml:"no-synced-lyrics"`
	LoadingLyrics  string `toml:"loading-lyrics"`
	ErrorMessage   string `toml:"error-message"`
	ServerOffline  string `toml:"server-offline"`
}

var GConfig struct {
	M sync.Mutex
	C Config
}

func readConfig() {
	ucd, _ := os.UserConfigDir()
	config, err := os.ReadFile(ucd + "/lrcsnc-lbl-client/config.toml")
	if errors.Is(err, os.ErrNotExist) {
		log.Fatal("ERROR: Config file does not exist. Ensure the config.toml file is in something like '$XDG_CONFIG_DIR/lrcsnc-lbl-client/config.toml'.")
	}

	if err := toml.Unmarshal(config, &GConfig.C); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			lines := strings.Join(strings.Split(decodeErr.String(), "\n"), "\n\t")
			log.Fatal("ERROR: Error parsing the config file: \n\t" + lines)
		}
	}

	validate()
}

func validate() {
	if GConfig.C.InStopAt >= GConfig.C.OutStartAt {
		log.Fatal("ERROR: in-stop-at should be less than out-start-at")
	}
}
