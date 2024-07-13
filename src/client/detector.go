package client

import (
	"fmt"
	"math"
)

type Detector struct {
	mean         float64
	std          float64
	samples      []float64
	samplesCount int
}

func NewDetector() *Detector {
	return &Detector{
		mean:    0,
		std:     1,
		samples: nil,
	}
}

func (d *Detector) Mean() float64 {
	return d.mean
}

func (d *Detector) Std() float64 {
	return d.std
}

func (d *Detector) Samples() []float64 {
	return d.samples
}

func (d *Detector) Update() error {
	if d.samplesCount <= 1 {
		return fmt.Errorf("samples size can't be <= 1, current is %d", len(d.samples))
	}

	// https://stats.stackexchange.com/questions/134476/how-to-estimate-mean-and-standard-deviation-of-a-normal-distribution-from-noisy
	d.mean = 0
	for i := 0; i < d.samplesCount; i++ {
		d.mean += d.samples[i]
	}
	d.mean /= float64(d.samplesCount)

	d.std = 0
	for i := 0; i < d.samplesCount; i++ {
		d.std += (d.samples[i] - d.mean) * (d.samples[i] - d.mean)
	}
	d.std = math.Sqrt(d.std / float64(d.samplesCount-1))

	return nil
}
