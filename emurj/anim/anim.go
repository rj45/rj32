package anim

import "time"

type Anim struct {
	frames KeyFrames
	T      float32
	Loop   bool
}

func (a *Anim) Next(V float32, dt time.Duration, easings ...Easing) {
	t := float32(dt.Seconds())
	if len(a.frames) > 0 {
		t += a.frames[len(a.frames)-1].T
	}
	a.frames = a.frames.Next(V, t, easings...)
}

func (a *Anim) Value() float32 {
	return a.frames.Value(a.T)
}

func (a *Anim) Advance(dt time.Duration) {
	a.T += float32(dt.Seconds())

	if a.Loop && len(a.frames) > 0 && a.T > a.frames[len(a.frames)-1].T {
		a.T -= a.frames[len(a.frames)-1].T
	}
}
