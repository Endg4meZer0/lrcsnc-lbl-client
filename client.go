package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/Endg4meZer0/go-mpris"
)

var conn net.Conn
var actualLyric string
var timeLeft float64
var rate float64 = 1

func initClient() {
	var err error
	conn, err = net.Dial(GConfig.C.Protocol, GConfig.C.ListenAt)
	if err != nil {
		log.Fatal("ERROR: Failed to connect to server: " + err.Error())
	}

	// Start the event reader
	go func() {
		reader := bufio.NewReader(conn)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("output/client", "Reading data from server returned an error: "+err.Error())
				return
			}

			var e Event
			err = json.Unmarshal([]byte(str), &e)
			if err != nil {
				log.Println("output/client", "Unexpected data received: "+str)
				overwrite(GConfig.C.Errors.ErrorMessage)
			} else {
				go receiveEvent(e)
			}
		}
	}()
}

func receiveEvent(e Event) {
	switch e.Type {
	case EventTypeServerClosed:
		conn.Close()
		overwrite(GConfig.C.Errors.ServerOffline)
		os.Exit(0)
	case EventTypeLyricsStateChanged:
		lyricsState = LyricsState(e.Data["State"].(float64))
		switch lyricsState {
		case LyricsStateInstrumental, LyricsStateSynced:
			text.status = StatusOK
		case LyricsStatePlain:
			text.status = StatusPlain
		case LyricsStateNotFound:
			text.status = StatusNotFound
		case LyricsStateLoading:
			text.status = StatusLoading
		case LyricsStateUnknown:
			text.status = StatusError
		}
	case EventTypePlayerChanged:
		playerName = e.Data["Name"].(string)
	case EventTypePlaybackStatusChanged:
		playbackStatus = mpris.PlaybackStatus(e.Data["PlaybackStatus"].(string))
		if playbackStatus == mpris.PlaybackStopped {
			overwrite(GConfig.C.Errors.NotPlaying)
		}
	case EventTypeOverwriteRequired:
		overwriteText = e.Data["Overwrite"].(string)
		if overwriteText != "" {
			overwrite(overwriteText)
		} else {
			text.index = len(actualLyric) - 1
			text.full = actualLyric
			text.inProgress = false
			write()
		}
	case EventTypeRateChanged:
		rate = e.Data["Rate"].(float64)
	case EventTypeActiveLyricChanged:
		lyric := e.Data["Lyric"].(map[string]any)["Text"].(string)
		actualLyric = lyric
		timeLeft = e.Data["TimeUntilEnd"].(float64)
		multiplier = int(e.Data["Multiplier"].(float64))
		if actualLyric == "" || e.Data["Index"].(float64) == -1 {
			text.m.Lock()
			text.full = ""
			text.inProgress = false
			text.index = 0
			text.m.Unlock()
			return
		}
		text.m.Lock()
		if overwriteText == "" && !e.Data["Resync"].(bool) {
			text.index = 0
			text.inProgress = true
			text.full = actualLyric
			go func() {
				startCounting()
				text.m.Unlock()
			}()
		} else if overwriteText == "" && e.Data["Resync"].(bool) {
			if actualLyric != text.full {
				text.index = 0
				text.inProgress = true
				text.full = actualLyric
				go func() {
					startCounting()
					text.m.Unlock()
				}()
			} else if text.inProgress && text.index > 3 {
				text.index = len(text.full)
				text.inProgress = false
				text.m.Unlock()
			}
		}
		go func() {
			if GConfig.C.OutStartAt == 100 {
				if overwriteText == "" && text.full == e.Data["Lyric"].(map[string]any)["Text"].(string) {
					text.full = ""
					text.index = 0
					text.inProgress = false
				}
				return
			}
			<-time.After(time.Duration((float64(GConfig.C.OutStartAt)/100*timeLeft)*1000/rate) * time.Millisecond)
			text.m.Lock()
			if overwriteText == "" && text.full == e.Data["Lyric"].(map[string]any)["Text"].(string) {
				text.inProgress = true
				startCounting()
			}
			text.m.Unlock()
		}()
	}
}
