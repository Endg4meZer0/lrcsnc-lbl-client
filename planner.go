package main

import (
	"context"
	"strings"
	"time"
)

type AnimationData struct {
	Value string
	Type  AnimationType
	Time  float64
	Force bool
}

type Planner struct {
	Lyric   *Animation
	Title   *Animation
	Artist  *Animation
	Artists *Animation
	Album   *Animation
}

type lyricQueueData struct {
	value string
	time  float64
}

var P Planner = Planner{
	Lyric:   NewStub(),
	Title:   NewStub(),
	Artist:  NewStub(),
	Artists: NewStub(),
	Album:   NewStub(),
}

func (p *Planner) RequestAnimation(ad AnimationData) {
	var chosenAnimation *Animation
	var alwaysSwitch bool
	var animationSpeed float64

	switch ad.Type {
	case AnimationTypeLyric:
		if GConfig.C.StopAppearAt == 0 {
			return
		}

		if ad.Value == "" && ad.Force {
			if p.Lyric.Flow() {
				l := p.Lyric.Source()
				if !p.Lyric.IsFinished {
					p.Lyric.Cancel()
				}
				<-p.Lyric.Finished
				close(p.Lyric.Finished)
				if GConfig.C.DisappearOnSwitch {
					ctx, canc := context.WithCancel(context.Background())
					p.Lyric = NewAnimation(ctx, canc, ad.Type, l, -1, AnimationFlowDisappear, 1)
					p.Lyric.Start()
				} else {
					wd.data["lyric"] = ""
					write()
				}
				return
			} else {
				return
			}
		}

		if ad.Force && !p.Lyric.IsFinished {
			p.Lyric.Cancel()
		}

		if ad.Time <= 0.33 || GConfig.C.StopAppearAt == 0 {
			wd.data["lyric"] = ad.Value
			write()
			if ad.Time <= 0.33 {
				return
			}
		}

		<-p.Lyric.Finished
		close(p.Lyric.Finished)
		ctx, canc := context.WithCancel(context.Background())
		p.Lyric = NewAnimation(ctx, canc, ad.Type, ad.Value, -1, AnimationFlowAppear, ad.Time)
		p.Lyric.Start()

		if GConfig.C.StartDisappearAt != 100 {
			go func() {
				l := ad.Value
				t := ((ad.Time-(float64(GConfig.C.StartDisappearAt)/100*ad.Time))*1000 - 100) / rate
				<-time.After(time.Duration(t) * time.Millisecond)
				if p.Lyric.Source() != l {
					return
				}
				<-p.Lyric.Finished
				close(p.Lyric.Finished)
				ctx, canc := context.WithCancel(context.Background())
				p.Lyric = NewAnimation(ctx, canc, ad.Type, l, -1, AnimationFlowDisappear, ad.Time)
				p.Lyric.Start()
			}()
		}

		return
	case AnimationTypeTitle:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeTitle)) {
			return
		}
		chosenAnimation = p.Title
		alwaysSwitch = GConfig.C.Info.AlwaysSwitch.Title
		animationSpeed = GConfig.C.Info.AnimationSpeed.Title
	case AnimationTypeArtist:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeArtist)) {
			return
		}
		chosenAnimation = p.Artist
		alwaysSwitch = GConfig.C.Info.AlwaysSwitch.Artist
		animationSpeed = GConfig.C.Info.AnimationSpeed.Artist
	case AnimationTypeArtists:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeArtists)) {
			return
		}
		chosenAnimation = p.Artists
		alwaysSwitch = GConfig.C.Info.AlwaysSwitch.Artists
		animationSpeed = GConfig.C.Info.AnimationSpeed.Artists
	case AnimationTypeAlbum:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeAlbum)) {
			return
		}
		chosenAnimation = p.Album
		alwaysSwitch = GConfig.C.Info.AlwaysSwitch.Album
		animationSpeed = GConfig.C.Info.AnimationSpeed.Album
	}

	if ad.Value == chosenAnimation.Source() && !alwaysSwitch {
		return
	}

	if chosenAnimation.Flow() {
		if !chosenAnimation.IsFinished {
			chosenAnimation.Cancel()
		}
		<-chosenAnimation.Finished
		close(chosenAnimation.Finished)
		ctx, canc := context.WithCancel(context.Background())
		var t float64
		if animationSpeed == 0 {
			t = 0
		} else {
			t = float64(len(chosenAnimation.Source())) / float64(animationSpeed)
		}
		chosenAnimation = NewAnimation(ctx, canc, ad.Type, chosenAnimation.Source(), -1, AnimationFlowDisappear, t)
		chosenAnimation.Start()
	}

	<-chosenAnimation.Finished
	close(chosenAnimation.Finished)
	ctx, canc := context.WithCancel(context.Background())
	var t float64
	if animationSpeed == 0 {
		t = 0
	} else {
		t = float64(len(chosenAnimation.Source())) / float64(animationSpeed)
	}
	chosenAnimation = NewAnimation(ctx, canc, ad.Type, ad.Value, -1, AnimationFlowAppear, t)
	chosenAnimation.Start()
}

func (*Planner) prepareLyric(lyric string) string {
	if len(lyric) == 0 {
		return lyric
	}

	switch GConfig.C.ForceCapitalization {
	case ForceCapitalizationTypeLower:
		lyric = strings.ToLower(lyric)
	case ForceCapitalizationTypeUpper:
		lyric = strings.ToUpper(lyric)
	}

	return lyric
}
