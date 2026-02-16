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
	Template            string                  `toml:"template"`
	MultiplierTemplate  string                  `toml:"multiplier"`
	ForceCapitalization ForceCapitalizationType `toml:"force-capitalization"`
	Protocol            string                  `toml:"protocol"`
	ListenAt            string                  `toml:"listen-at"`
	StopAppearAt        int                     `toml:"stop-appear-at"`
	StartDisappearAt    int                     `toml:"start-disappear-at"`
	DisappearOnSwitch   bool                    `toml:"disappear-on-switch"`
	Info                InfoConfig              `toml:"info"`
	Instrumental        InstrumentalConfig      `toml:"instrumental"`
	Errors              ErrorsConfig            `toml:"errors"`
}

type InfoConfig struct {
	AnimationSpeed InfoAnimationSpeedConfig `toml:"animation-speed"`
	AlwaysSwitch   InfoAlwaysSwitchConfig   `toml:"always-switch"`
}

type InfoAnimationSpeedConfig struct {
	Title   float64 `toml:"title"`
	Artist  float64 `toml:"artist"`
	Artists float64 `toml:"artists"`
	Album   float64 `toml:"album"`
}

type InfoAlwaysSwitchConfig struct {
	Title   bool `toml:"title"`
	Artist  bool `toml:"artist"`
	Artists bool `toml:"artists"`
	Album   bool `toml:"album"`
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
	ServerError    string `toml:"server-error"`
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
	if GConfig.C.StopAppearAt >= GConfig.C.StartDisappearAt {
		log.Fatal("ERROR: stop-appear-at should be less than start-disappear-at")
	}
	if GConfig.C.ForceCapitalization != ForceCapitalizationTypeNone &&
		GConfig.C.ForceCapitalization != ForceCapitalizationTypeUpper &&
		GConfig.C.ForceCapitalization != ForceCapitalizationTypeLower {
		log.Fatal("ERROR: force-capitalization's value is invalid")
	}
}

func (c *Config) TemplateHasKey(key string) bool {
	return strings.Contains(c.Template, "%"+key+"%")
}

// types

type ForceCapitalizationType string

const (
	ForceCapitalizationTypeNone  ForceCapitalizationType = "none"
	ForceCapitalizationTypeUpper ForceCapitalizationType = "uppercase"
	ForceCapitalizationTypeLower ForceCapitalizationType = "lowercase"
)
