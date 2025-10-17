package main

import "strings"

type Lyric struct {
	Timing float64
	Text   string
}

type LyricsState byte

const (
	LyricsStateSynced       LyricsState = 0
	LyricsStatePlain        LyricsState = 1
	LyricsStateInstrumental LyricsState = 2
	LyricsStateNotFound     LyricsState = 3
	LyricsStateLoading      LyricsState = 4
	LyricsStateUnknown      LyricsState = 5
)

var lyricsStateStrings = map[LyricsState]string{
	LyricsStateSynced:       "synced",
	LyricsStatePlain:        "plain",
	LyricsStateInstrumental: "instrumental",
	LyricsStateNotFound:     "not-found",
	LyricsStateLoading:      "loading",
	LyricsStateUnknown:      "unknown",
}

func (l LyricsState) String() string {
	return lyricsStateStrings[l]
}

type Status byte

const (
	StatusOK Status = iota
	StatusPlain
	StatusNotFound
	StatusLoading
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusPlain:
		return GConfig.C.Errors.NoSyncedLyrics
	case StatusNotFound:
		return GConfig.C.Errors.NoLyrics
	case StatusLoading:
		return GConfig.C.Errors.LoadingLyrics
	case StatusError:
		return GConfig.C.Errors.ErrorMessage
	default:
		return ""
	}
}

// DynamicReplacer is a utility used to improve template functionality
// of output/client module. It is a very simplified middle-ground between
// strings.Replacer and a proper text/template.
//
// It is hardcoded to use curly brackets as start and end delimeters for keys.
type DynamicReplacer struct {
	m map[string]func() string
}

const keyDelim = '%'

func NewDynamicReplacer(_m map[string]func() string) *DynamicReplacer {
	return &DynamicReplacer{
		m: _m,
	}
}

func (dr *DynamicReplacer) Replace(template string) string {
	var result strings.Builder

	// An optimization to reduce the number of allocations.
	result.Grow(len(template) * 2)

	i := 0

	for i < len(template) {
		startIndex := strings.IndexRune(template[i:], keyDelim)
		if startIndex == -1 {
			result.WriteString(template[i:])
			break
		}
		startIndex += i
		result.WriteString(template[i:startIndex])

		endIndex := strings.IndexRune(template[startIndex+1:], keyDelim)
		if endIndex == -1 {
			result.WriteString(template[startIndex:])
			break
		}
		endIndex += startIndex + 1
		key := template[startIndex+1 : endIndex]

		if fn, ok := dr.m[key]; ok {
			res := fn()
			result.WriteString(res)
		} else {
			result.WriteString(template[startIndex : endIndex+1])
		}
		i = endIndex + 1
	}

	return result.String()
}
