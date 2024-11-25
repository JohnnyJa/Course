package Interface

type IProcess interface {
	RunToCurrentTime()

	UpdateCurrentTime(float64)

	GetActivationTime() float64
}
