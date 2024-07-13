package client

import (
	"fmt"
	"log"
	"math"
	"sync"
	frequency "team00/generated"
)

var frequenciesBuffer = sync.Pool{
	New: func() any { return make([]float64, DetectSamplesOptimalSize) },
}

const (
	DetectSamplesOptimalSize int = 200
)

type Client struct {
	Detector         Detector
	IsDetectFinished bool
}

func NewClient() *Client {
	return &Client{
		Detector:         *NewDetector(),
		IsDetectFinished: false,
	}
}

func (c *Client) DetectAnomaly(freq, k float64) (bool, float64, error) {
	if !c.IsDetectFinished {
		return false, 0, fmt.Errorf("detect phase is not finished")
	}

	diff := math.Abs(math.Abs(freq) - math.Abs(c.Detector.mean))
	return diff > float64(k)*c.Detector.std, diff, nil
}

func (c *Client) Calibrate(stream frequency.FrequencyRandomizer_SpawnFrequenciesClient, calibrateOn int) error {
	if c.IsDetectFinished {
		// avoid repeated calibration
		return nil
	}

	c.Detector.samples = frequenciesBuffer.Get().([]float64)
	defer frequenciesBuffer.Put(&c.Detector.samples)

	for i := 0; i < calibrateOn; i++ {
		message, err := stream.Recv()
		if err != nil {
			return err
		}

		c.Detector.samples[i] = message.GetFrequency()
		c.Detector.samplesCount = i + 1
		if err := c.Detector.Update(); err == nil {
			log.Printf("calibrated on %d samples: mean=%f; std=%f\n", c.Detector.samplesCount, c.Detector.Mean(), c.Detector.Std())
		}
	}

	log.Printf("calibration finished: mean=%f; std=%f\n", c.Detector.Mean(), c.Detector.Std())
	c.IsDetectFinished = true
	return nil
}
