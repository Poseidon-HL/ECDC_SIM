package util

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Weibull struct {
	shape    float64
	scale    float64
	location float64
}

func NewWeibull(shape, scale, location float64) *Weibull {
	return &Weibull{
		shape:    shape,
		scale:    scale,
		location: location,
	}
}

func (w *Weibull) pdf(x float64) float64 {
	if x < 0 || x < w.location {
		return 0
	}
	factorA := w.shape / w.scale
	factorB := math.Pow((x-w.location)/w.scale, w.shape-1)
	factorC := math.Exp(-math.Pow((x-w.location)/w.scale, w.shape))
	return factorA * factorB * factorC
}

func (w *Weibull) cdf(x float64) float64 {
	if x < w.location {
		return 0
	}
	return 1 - math.Exp(-math.Pow((x-w.location)/w.scale, w.shape))
}

func (w *Weibull) HazardRate(x float64) float64 {
	if x < w.location {
		return 0
	}
	if w.shape == 1 {
		return 1 / w.scale
	}
	return math.Abs(w.pdf(x) / (1 - w.cdf(x)))
}

func (w *Weibull) Draw() float64 {
	u := 1 - rand.Float64()
	return w.scale*math.Pow(-math.Log(u), 1/w.shape) + w.location
}
