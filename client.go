package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Endg4meZer0/go-mpris"
)

var conn net.Conn
var rate float64 = 1

func initClient() {
	var err error
	conn, err = net.Dial(GConfig.C.Protocol, GConfig.C.ListenAt)
	if err != nil {
		log.Println("ERROR: Failed to connect to server: " + err.Error())
		go receiveEvent(Event{
			Type: EventTypeServerError,
			Data: map[string]any{},
		})
		return
	}

	// Start the event reader
	go func() {
		reader := bufio.NewReader(conn)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				log.Println("output/client", "Reading data from server returned an error: "+err.Error())
				receiveEvent(Event{
					Type: EventTypeServerError,
					Data: map[string]any{},
				})
				return
			}

			var e Event
			err = json.Unmarshal([]byte(str), &e)
			if err != nil {
				log.Println("output/client", "Unexpected data received: "+str)
				message(GConfig.C.Errors.ErrorMessage)
			} else {
				receiveEvent(e)
			}
		}
	}()
}

func receiveEvent(e Event) {
	switch e.Type {
	case EventTypeServerClosed:
		conn.Close()
		message(GConfig.C.Errors.ServerOffline)
		os.Exit(0)
	case EventTypeServerError:
		_ = conn.Close()
		message(GConfig.C.Errors.ServerError)
		go func() {
			<-time.After(10 * time.Second)
			initClient()
		}()
	case EventTypeLyricsStateChanged:
		lyricsState = LyricsState(e.Data["State"].(float64))
		switch lyricsState {
		case LyricsStateInstrumental, LyricsStateSynced:
			wd.status = StatusOK
		case LyricsStatePlain:
			wd.status = StatusPlain
		case LyricsStateNotFound:
			wd.status = StatusNotFound
		case LyricsStateLoading:
			wd.status = StatusLoading
		case LyricsStateUnknown:
			wd.status = StatusError
		}
	case EventTypeSongChanged:
		title := e.Data["Title"].(string)
		album := e.Data["Album"].(string)
		artists := make([]string, len(e.Data["Artists"].([]any)))
		for i, a := range e.Data["Artists"].([]any) {
			artists[i] = a.(string)
		}
		go P.RequestAnimation(AnimationData{
			Value: title,
			Type:  AnimationTypeTitle,
			Time:  -1,
			Force: true,
		})
		go P.RequestAnimation(AnimationData{
			Value: album,
			Type:  AnimationTypeAlbum,
			Time:  -1,
			Force: true,
		})
		go P.RequestAnimation(AnimationData{
			Value: strings.Join(artists, ", "),
			Type:  AnimationTypeArtists,
			Time:  -1,
			Force: true,
		})
		if len(artists) > 0 {
			go P.RequestAnimation(AnimationData{
				Value: artists[0],
				Type:  AnimationTypeArtist,
				Time:  -1,
				Force: true,
			})
		}
	case EventTypePlayerChanged:
		playerName = e.Data["Name"].(string)
		if playerName == "" {
			wd.status = StatusNotPlaying
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeTitle,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeArtist,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeArtists,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeAlbum,
				Time:  -1,
				Force: true,
			})
		}
	case EventTypePlaybackStatusChanged:
		playbackStatus = mpris.PlaybackStatus(e.Data["PlaybackStatus"].(string))
		if playbackStatus == mpris.PlaybackStopped {
			wd.status = StatusNotPlaying
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeTitle,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeArtist,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeArtists,
				Time:  -1,
				Force: true,
			})
			go P.RequestAnimation(AnimationData{
				Value: "",
				Type:  AnimationTypeAlbum,
				Time:  -1,
				Force: true,
			})
		}
	case EventTypeOverwriteRequired:
		overwriteText = e.Data["Overwrite"].(string)
		overwrite(overwriteText)
	case EventTypeRateChanged:
		rate = e.Data["Rate"].(float64)
	case EventTypeActiveLyricChanged:
		lyric := e.Data["Lyric"].(map[string]any)["Text"].(string)
		duration := e.Data["TimeUntilEnd"].(float64)
		multiplier = int(e.Data["Multiplier"].(float64))

		go P.RequestAnimation(AnimationData{
			Value: lyric,
			Type:  AnimationTypeLyric,
			Time:  duration,
			Force: e.Data["Resync"].(bool),
		})

	}
}
