package rog

import (
	"time"
)

type stats struct {
	Then, Now   time.Time
	Elapsed, Dt float64
	Frames, Fps int64
}

func (t *stats) Update() {
	now := time.Now()
	t.Then = t.Now
	t.Now = now
	t.Dt = t.Now.Sub(t.Then).Seconds()
	t.Elapsed += t.Dt
	t.Frames += 1
	if t.Elapsed >= 1 {
		t.Fps = t.Frames
		t.Frames = 0
		t.Elapsed -= t.Elapsed
	}
}
