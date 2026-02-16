package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Endg4meZer0/go-mpris"
)

type writerFields struct {
	lyric   string
	title   string
	artist  string
	artists string
	album   string
	instr   string
}

type writerData struct {
	overwrite    string
	message      string
	messageShown bool
	status       Status
	fields       writerFields
}

var wd writerData
var write func() = writer()

var multiplier int
var overwriteText string

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

				wd.fields.instr = strings.Repeat(GConfig.C.Instrumental.Symbol, i)
				write()

				i++
				if i > j {
					i = 1
				}
			}
		}()
	} else {
		instrTicker.Stop()
		wd.fields.instr = ""
	}
}

func message(txt string) {
	wd.message = txt
	write()
	go func() {
		<-time.After(5 * time.Second)
		wd.message = ""
		wd.messageShown = false
		write()
	}()
}

func overwrite(txt string) {
	wd.overwrite = txt
	write()
}

func writer() func() {
	var dynrepl *DynamicReplacer = NewDynamicReplacer(
		map[string]func() string{
			"text": func() string {
				if wd.overwrite != "" {
					return strings.ReplaceAll(wd.overwrite, "\"", "\\\"")
				}

				if wd.status != StatusOK {
					return wd.status.String()
				}

				if wd.message != "" {
					return strings.ReplaceAll(wd.message, "\"", "\\\"")
				}

				return strings.ReplaceAll(wd.fields.lyric, "\"", "\\\"")
			},
			"title": func() string {
				return strings.ReplaceAll(wd.fields.title, "\"", "\\\"")
			},
			"artist": func() string {
				return strings.ReplaceAll(wd.fields.artist, "\"", "\\\"")
			},
			"artists": func() string {
				return strings.ReplaceAll(wd.fields.artists, "\"", "\\\"")
			},
			"album": func() string {
				return strings.ReplaceAll(wd.fields.album, "\"", "\\\"")
			},
			"instr": func() string {
				return wd.fields.instr
			},
			"player": func() string {
				return playerName
			},
			"multiplier": func() string {
				if multiplier <= 1 {
					return ""
				}

				return strings.ReplaceAll(GConfig.C.MultiplierTemplate, "%value%", fmt.Sprint(multiplier))
			},
		},
	)

	return func() {
		fmt.Println(dynrepl.Replace(GConfig.C.Template))
	}
}
