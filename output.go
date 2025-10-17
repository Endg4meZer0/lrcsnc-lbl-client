package main

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/Endg4meZer0/go-mpris"
)

type animated struct {
	m          sync.Mutex
	status     Status
	index      int
	full       string
	inProgress bool
}

var text animated
var instrumental string
var multiplier int
var overwriteText string
var dynrepl *DynamicReplacer = NewDynamicReplacer(
	map[string]func() string{
		"text": func() string {
			out := ""

			if text.status != StatusOK {
				return text.status.String()
			}

			if !text.inProgress && text.index == 0 {
				return ""
			}

			if text.inProgress {
				out = string([]rune(text.full)[:text.index])
			} else {
				out = text.full
			}

			out = strings.ReplaceAll(out, "\"", "\\\"")

			if GConfig.C.Multiplier.Enabled && GConfig.C.Multiplier.AddTo == "text" {
				return addMultiplier(out)
			} else {
				return out
			}
		},
		"instr": func() string {
			out := instrumental
			if GConfig.C.Multiplier.Enabled && GConfig.C.Multiplier.AddTo == "instr" {
				return addMultiplier(out)
			} else {
				return out
			}
		},
		"player": func() string {
			return playerName
		},
	},
)

var instrTicker = time.NewTicker(5 * time.Minute)

func initOutput() {
	if GConfig.C.Instrumental.Enabled {
		instrTicker.Reset(time.Duration(GConfig.C.Instrumental.Interval*1000) * time.Millisecond)
		go func() {
			i := 1
			j := GConfig.C.Instrumental.MaxSymbols
			for {
				<-instrTicker.C
				if playbackStatus != mpris.PlaybackPlaying {
					continue
				}

				instrumental = strings.Repeat(GConfig.C.Instrumental.Symbol, i)
				write()

				i++
				if i > j {
					i = 1
				}
			}
		}()
	} else {
		instrTicker.Stop()
		instrumental = ""
	}
}

func overwrite(s string) {
	text.full = s
	text.inProgress = false
}

func write() {
	fmt.Println(dynrepl.Replace(GConfig.C.Template))
}

func addMultiplier(s string) string {
	if multiplier < 2 {
		return s
	}

	if GConfig.C.Multiplier.ToLeft {
		return strings.ReplaceAll(GConfig.C.Multiplier.Format, "%value%",
			fmt.Sprintf("%v", multiplier)) + s
	} else {
		return s + strings.ReplaceAll(GConfig.C.Multiplier.Format, "%value%",
			fmt.Sprintf("%v", multiplier))
	}
}

func startCounting() {
	txtRunes := []rune(text.full)

	if text.index < len(txtRunes) {
		if len(txtRunes) == 1 {
			text.index = 1
			text.inProgress = false
			return
		}
		timeBetweenLettersMs := int(math.Round((timeLeft*float64(GConfig.C.InStopAt)/100)*1000/rate)) / (len(txtRunes) - text.index - 1)
		if timeBetweenLettersMs <= 0 {
			text.index = len(txtRunes)
			text.inProgress = false
			return
		}
		ticker := time.NewTicker(time.Duration(timeBetweenLettersMs) * time.Millisecond)
		for range ticker.C {
			text.index++
			if text.index > len(txtRunes) || !text.inProgress {
				text.index = len(txtRunes)
				text.inProgress = false
				break
			}
			write()
		}
		ticker.Stop()
	} else {
		if len(txtRunes) == 1 {
			text.index = 0
			text.inProgress = false
			return
		}
		timeBetweenLettersMs := int(math.Round((timeLeft-timeLeft*float64(GConfig.C.OutStartAt)/100)*1000/rate)) / len(txtRunes)
		if timeBetweenLettersMs <= 0 {
			text.index = 0
			text.inProgress = false
			return
		}
		ticker := time.NewTicker(time.Duration(timeBetweenLettersMs) * time.Millisecond)
		for range ticker.C {
			text.index--
			if text.index < 0 || !text.inProgress {
				text.index = 0
				text.inProgress = false
				break
			}
			write()
		}
		ticker.Stop()
	}
}
