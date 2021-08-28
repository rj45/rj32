package anim

func Clamp01(t float32) float32 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func Lerp(start, end, pct float32) float32 {
	return start + ((end - start) * Clamp01(pct))
}

// TODO: make this work on an array of floats
// so multiple properties can be animated at once
// - anim takes list of property #s, list of values, duration

type KeyFrame struct {
	V    float32
	T    float32
	Ease Easing
}

func (kf *KeyFrame) Lerp(next *KeyFrame, t float32) float32 {
	tp := t - kf.T
	dt := next.T - kf.T

	if dt == 0 {
		dt = 1
	}

	pct := tp / dt

	return Lerp(kf.V, next.V, kf.Ease(pct))
}

type KeyFrames []KeyFrame

func (a KeyFrames) Value(t float32) float32 {
	frame := 0
	for t > a[frame].T && frame < len(a) {
		frame++
	}

	if a[frame].T > t && frame > 0 {
		frame--
	}

	if frame == len(a)-1 || len(a) < 1 {
		return a[frame].V
	}

	return a[frame].Lerp(&a[frame+1], t)
}

func (kf KeyFrames) Next(V, T float32, easings ...Easing) KeyFrames {
	var easing Easing

	switch len(easings) {
	case 0:
		easing = None
	case 1:
		easing = easings[0]
	default:
		panic("todo: only 0 or 1 easing is implemented")
	}

	return append(kf, KeyFrame{V: V, T: T, Ease: easing})
}
