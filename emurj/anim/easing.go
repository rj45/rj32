package anim

type Easing func(t float32) float32

func EaseIn(t float32) float32 {
	return t * t
}

func Flip(t float32) float32 {
	return 1 - t
}

func EaseOut(t float32) float32 {
	return Flip(EaseIn(Flip(t)))
}

func EaseInOut(t float32) float32 {
	return Lerp(EaseIn(t), EaseOut(t), t)
}

func None(t float32) float32 {
	return t
}
