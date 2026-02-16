package main

import (
	"context"
	"strings"
	"sync"
	"time"
)

type AnimationData struct {
	Value string
	Type  AnimationType
	Time  float64
	Force bool
}

type PlannerAnim struct {
	M sync.Mutex
	A *Animation
}

type Planner struct {
	Lyric   PlannerAnim
	Title   PlannerAnim
	Artist  PlannerAnim
	Artists PlannerAnim
	Album   PlannerAnim
}

type lyricQueueData struct {
	value string
	time  float64
}

var P Planner = Planner{
	Lyric:   PlannerAnim{A: NewStub()},
	Title:   PlannerAnim{A: NewStub()},
	Artist:  PlannerAnim{A: NewStub()},
	Artists: PlannerAnim{A: NewStub()},
	Album:   PlannerAnim{A: NewStub()},
}

func (p *Planner) RequestAnimation(ad AnimationData) {
	switch ad.Type {
	case AnimationTypeLyric:
		p.Lyric.M.Lock()
		defer p.Lyric.M.Unlock()
		if ad.Value == "" && ad.Force {
			if p.Lyric.A.Flow() {
				l := p.Lyric.A.Source()
				if !p.Lyric.A.IsFinished {
					p.Lyric.A.Cancel()
				}
				<-p.Lyric.A.Finished
				close(p.Lyric.A.Finished)
				if GConfig.C.DisappearOnSwitch && ad.Time > 0.33 {
					ctx, canc := context.WithCancel(context.Background())
					t := float64(1) / (1 - float64(GConfig.C.StartDisappearAt)/100)
					if t < ad.Time {
						t = ad.Time - 0.1
					}
					p.Lyric.A = NewAnimation(ctx, canc, ad.Type, l, -1, AnimationFlowDisappear, t)
					go p.Lyric.A.Start()
				} else {
					wd.fields.lyric = ""
					write()
					p.Lyric.A = NewStub()
				}
				return
			} else {
				return
			}
		}

		if ad.Force && !p.Lyric.A.IsFinished {
			p.Lyric.A.Cancel()
		}

		<-p.Lyric.A.Finished
		close(p.Lyric.A.Finished)

		ctx, canc := context.WithCancel(context.Background())
		p.Lyric.A = NewAnimation(ctx, canc, ad.Type, ad.Value, -1, AnimationFlowAppear, ad.Time)
		go p.Lyric.A.Start()

		if GConfig.C.StartDisappearAt != 100 {
			go func() {
				l := ad.Value
				<-time.After(time.Duration(((float64(GConfig.C.StartDisappearAt)/100*ad.Time)*1000-100)/rate) * time.Millisecond)
				if p.Lyric.A.Source() != l || !p.Lyric.A.Flow() {
					return
				}
				<-p.Lyric.A.Finished
				close(p.Lyric.A.Finished)
				ctx, canc := context.WithCancel(context.Background())
				p.Lyric.A = NewAnimation(ctx, canc, ad.Type, l, -1, AnimationFlowDisappear, ad.Time)
				go p.Lyric.A.Start()
			}()
		}

		return
	case AnimationTypeTitle:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeTitle)) {
			return
		}
		p.Title.M.Lock()
		defer p.Title.M.Unlock()
		p.nonLyricAnimation(&p.Title.A, ad, GConfig.C.Info.AlwaysSwitch.Title, GConfig.C.Info.AnimationSpeed.Title)
	case AnimationTypeArtist:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeArtist)) {
			return
		}
		p.Artist.M.Lock()
		defer p.Artist.M.Unlock()
		p.nonLyricAnimation(&p.Artist.A, ad, GConfig.C.Info.AlwaysSwitch.Artist, GConfig.C.Info.AnimationSpeed.Artist)
	case AnimationTypeArtists:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeArtists)) {
			return
		}
		p.Artists.M.Lock()
		defer p.Artists.M.Unlock()
		p.nonLyricAnimation(&p.Artists.A, ad, GConfig.C.Info.AlwaysSwitch.Artists, GConfig.C.Info.AnimationSpeed.Artists)
	case AnimationTypeAlbum:
		if !GConfig.C.TemplateHasKey(string(AnimationTypeAlbum)) {
			return
		}
		p.Album.M.Lock()
		defer p.Album.M.Unlock()
		p.nonLyricAnimation(&p.Album.A, ad, GConfig.C.Info.AlwaysSwitch.Album, GConfig.C.Info.AnimationSpeed.Album)
	}
}

// this abomination is just insane and yet it works perfectly as I want it to
func (*Planner) nonLyricAnimation(an **Animation, ad AnimationData, alwaysSwitch bool, animationSpeed float64) {
	if ad.Value == (*an).Source() && !alwaysSwitch {
		return
	}

	if (*an).Flow() {
		if !(*an).IsFinished {
			(*an).Cancel()
		}
		<-(*an).Finished
		close((*an).Finished)
		ctx, canc := context.WithCancel(context.Background())
		var t float64
		if animationSpeed == 0 {
			t = 0
		} else {
			t = float64(len([]rune((*an).Source()))) / float64(animationSpeed)
		}
		*an = NewAnimation(ctx, canc, ad.Type, (*an).Source(), -1, AnimationFlowDisappear, t)
		go (*an).Start()
	}

	<-(*an).Finished
	close((*an).Finished)
	ctx, canc := context.WithCancel(context.Background())
	var t float64
	if animationSpeed == 0 {
		t = 0
	} else {
		t = float64(len([]rune(ad.Value))) / float64(animationSpeed)
	}
	*an = NewAnimation(ctx, canc, ad.Type, ad.Value, -1, AnimationFlowAppear, t)
	go (*an).Start()
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
