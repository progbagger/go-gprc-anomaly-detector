package anomalies

type Detector interface {
	Mean() float64
	Std() float64
	Samples() []float64

	Update() error
}
