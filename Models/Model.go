package Models

import (
	"Course/Interface"
	"log"
)

type Model struct {
	simulationTime float64
	elements       []Interface.IProcess
	CurrentTime    float64

	FinishCycle func()
}

func NewModel(simulationTime float64, elements []Interface.IProcess) *Model {
	return &Model{
		simulationTime: simulationTime,
		elements:       elements,
		CurrentTime:    0.0,
	}
}

func (m *Model) Simulate() {
	for m.CurrentTime < m.simulationTime {
		m.CurrentTime = m.FindNextActivationTime()
		m.RunToCurrentTime(m.CurrentTime)

		m.FinishCycle()
	}

}

func (m *Model) FindNextActivationTime() float64 {
	minTime := m.simulationTime
	for _, element := range m.elements {
		if element.GetActivationTime() < minTime {
			minTime = element.GetActivationTime()
		}
	}
	return minTime
}

func (m *Model) RunToCurrentTime(time float64) {
	for _, element := range m.elements {
		element.UpdateCurrentTime(time)
	}
	log.Default().Printf("----Current time: %f----", time)
	for _, element := range m.elements {
		element.RunToCurrentTime()
	}
}
