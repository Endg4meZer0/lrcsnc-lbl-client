package main

import (
	"github.com/Endg4meZer0/go-mpris"
)

type Event struct {
	Type EventType
	Data map[string]any
}

// EventType represents the type of the event (received or sent)
type EventType string

const (
	EventTypeActiveLyricChanged    EventType = "ActiveLyricChanged"
	EventTypeSongChanged           EventType = "SongChanged"
	EventTypePlayerChanged         EventType = "PlayerChanged"
	EventTypePlaybackStatusChanged EventType = "PlaybackStatusChanged"
	EventTypeRateChanged           EventType = "RateChanged"
	EventTypeLyricsStateChanged    EventType = "LyricsStateChanged"
	EventTypeLyricsChanged         EventType = "LyricsChanged"
	EventTypeOverwriteRequired     EventType = "OverwriteRequired"
	EventTypeServerClosed          EventType = "ServerClosed"
	EventTypeConfigReloaded        EventType = "ConfigReloaded" // only for client
)

// These types are mostly only for references
// 'cause I do forget stuff I wrote in lrcsnc
// sometimes

type EventTypeActiveLyricChangedData struct {
	Index        int
	Lyric        Lyric
	Multiplier   int
	TimeUntilEnd float64
	Resync       bool
}

type EventTypeSongChangedData struct {
	Title    string
	Artists  []string
	Album    string
	Duration float64
}

type EventTypePlayerChangedData struct {
	Name string
}

type EventTypePlaybackStatusChangedData struct {
	PlaybackStatus mpris.PlaybackStatus
}

type EventTypeRateChangedData struct {
	Rate float64
}

type EventTypeLyricsStateChangedData struct {
	State LyricsState
}

type EventTypeLyricsChangedData struct {
	Lyrics []Lyric
}

type EventTypeOverwriteRequiredData struct {
	Overwrite string
}

type EventTypeServerClosedData struct{}
