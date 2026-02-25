package main

import (
	"context"
	"math"

	"time"
)

type AnimationType string

const (
	AnimationTypeLyric   AnimationType = "lyric"
	AnimationTypeTitle   AnimationType = "title"
	AnimationTypeArtist  AnimationType = "artist"
	AnimationTypeArtists AnimationType = "artists"
	AnimationTypeAlbum   AnimationType = "album"
)

type AnimationFlow bool

const (
	AnimationFlowAppear    AnimationFlow = true
	AnimationFlowDisappear               = false
)

// Animation is intended to be used ONCE.
// After an animation is done, you need to make a new one.
type Animation struct {
	ctx context.Context

	key      AnimationType
	source   string
	index    int
	flow     AnimationFlow
	duration float64

	Finished   chan bool
	IsFinished bool
	Cancel     func()
}

func NewStub() *Animation {
	an := &Animation{}
	an.Finished = make(chan bool, 1)
	an.Finished <- true
	an.IsFinished = true

	return an
}

func NewAnimation(ctx context.Context, canc func(), key AnimationType, source string, startIndex int, flow AnimationFlow, duration float64) *Animation {
	an := &Animation{}
	an.ctx = ctx
	an.Cancel = canc
	an.key = key
	an.source = source
	an.index = 0
	if startIndex != -1 {
		an.index = startIndex
	} else if !flow {
		an.index = len([]rune(source))
	}
	an.flow = flow
	an.duration = duration

	an.Finished = make(chan bool, 1)

	return an
}

func (an *Animation) Flow() AnimationFlow {
	return an.flow
}

func (an *Animation) Source() string {
	return an.source
}

func (an *Animation) Index() int {
	return an.index
}

func (an *Animation) Start() {
	runes := []rune(an.source)
	stopAppearAt := GConfig.C.StopAppearAt
	startDisappearAt := GConfig.C.StartDisappearAt
	if an.key != AnimationTypeLyric {
		stopAppearAt = 100
		startDisappearAt = 0
	}

	if an.flow {
		if stopAppearAt == 0 || an.duration == 0 || len(runes) == 0 {
			an.index = len(runes)
			an.writetxt()
			an.Finished <- true
			an.IsFinished = true
			return
		}

		// count the time between letters
		timeBetweenLettersMs := int(math.Round((an.duration*float64(stopAppearAt)/100)*1000/rate)) / len(runes)
		if timeBetweenLettersMs <= 0 {
			an.index = len(runes)
			an.writetxt()
			an.Finished <- true
			an.IsFinished = true
			return
		}

		an.index++
		an.writetxt()
		if an.index > len(runes) {
			an.Finished <- true
			an.IsFinished = true
			return
		}

		ticker := time.NewTicker(time.Duration(timeBetweenLettersMs) * time.Millisecond)
		defer ticker.Stop()
		for !an.IsFinished {
			select {
			case <-ticker.C:
				if an.IsFinished {
					break
				}

				an.index++
				if an.index > len(runes) {
					an.index = len(runes)
					an.writetxt()
					an.Finished <- true
					an.IsFinished = true
					break
				}
				an.writetxt()
			case <-an.ctx.Done():
				an.index = len(runes)
				an.writetxt()
				an.Finished <- true
				an.IsFinished = true
			}
		}
	} else {
		if startDisappearAt == 100 || an.duration == 0 || len(runes) == 0 {
			an.index = 0
			an.writetxt()
			an.Finished <- true
			an.IsFinished = true
			return
		}

		// count time between letters
		timeBetweenLettersMs := int(math.Round(((an.duration-0.1)-(an.duration-0.1)*float64(startDisappearAt)/100)*1000/rate)) / len(runes)
		if timeBetweenLettersMs <= 0 {
			an.index = 0
			an.writetxt()
			an.Finished <- true
			an.IsFinished = true
			return
		}

		an.index--
		an.writetxt()
		if an.index < 1 {
			an.Finished <- true
			an.IsFinished = true
			return
		}

		ticker := time.NewTicker(time.Duration(timeBetweenLettersMs) * time.Millisecond)
		defer ticker.Stop()
		for !an.IsFinished {
			select {
			case <-ticker.C:
				if an.IsFinished {
					break
				}

				an.index--
				if an.index < 0 {
					an.index = 0
					an.writetxt()
					an.Finished <- true
					an.IsFinished = true
					break
				}
				an.writetxt()
			case <-an.ctx.Done():
				an.index = 0
				an.writetxt()
				an.Finished <- true
				an.IsFinished = true
			}
		}
	}
}

func (an *Animation) writetxt() {
	switch an.key {
	case AnimationTypeLyric:
		wd.fields.lyric = string([]rune(an.source)[:an.index])
	case AnimationTypeTitle:
		wd.fields.title = string([]rune(an.source)[:an.index])
	case AnimationTypeArtist:
		wd.fields.artist = string([]rune(an.source)[:an.index])
	case AnimationTypeArtists:
		wd.fields.artists = string([]rune(an.source)[:an.index])
	case AnimationTypeAlbum:
		wd.fields.album = string([]rune(an.source)[:an.index])
	}
	write()
}
