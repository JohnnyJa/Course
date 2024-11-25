package Process

import (
	"Course/Marker"
	"log"
	"math"
)

type Process struct {
	Name               string
	CurrentTime        float64
	nextActivationTime float64

	Markers []*Marker.Marker

	GetDelay   func() float64
	EndProcess func()
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) RunToCurrentTime() {

	if p.GetActivationTime() <= p.CurrentTime {
		m := p.GetFinishedMarkers()
		arr := make([]int, len(m))
		for i := 0; i < len(m); i++ {
			arr[i] = m[i].Id
		}

		log.Default().Printf("Process %s end processing markers %d at %f", p.Name, arr, p.CurrentTime)

		p.EndProcess()

	} else {
		if len(p.Markers) > 0 {
			arr := make([]int, len(p.Markers))
			for i := 0; i < len(p.Markers); i++ {
				arr[i] = p.Markers[i].Id
			}
			log.Default().Printf("Process %s run processing marker %v at %f", p.Name, arr, p.CurrentTime)
		} else {
			log.Default().Printf("Process %s run do nothing at %f", p.Name, p.CurrentTime)
		}
	}
}

func (p *Process) GetActivationTime() float64 {
	minTime := math.MaxFloat64
	for i := 0; i < len(p.Markers); i++ {
		if p.Markers[i].NextActivationTime < minTime {
			minTime = p.Markers[i].NextActivationTime
		}
	}

	return minTime
}

//func (p *Process) TryTakeMarker(marker *Marker.Marker) bool {
//	if pr, ok := p.system.TryStartAtFreeProcessor(); ok {
//		p.StartProcess()
//
//		marker.NextActivationTime = p.GetDelay() + p.CurrentTime
//		marker.ProcessorsId = append(marker.ProcessorsId, pr)
//
//		p.Markers = append(p.Markers, marker)
//		return true
//	}
//
//	return false
//}

func (p *Process) TakeMarker(marker *Marker.Marker) {
	//p.StartProcess()
	log.Default().Printf("Process %s take marker %d at %f", p.Name, marker.Id, p.CurrentTime)

	marker.NextActivationTime = p.GetDelay() + p.CurrentTime
	p.Markers = append(p.Markers, marker)
}

func (p *Process) GetFinishedMarkers() []*Marker.Marker {
	result := make([]*Marker.Marker, 0)
	for i := 0; i < len(p.Markers); i++ {
		if p.Markers[i].NextActivationTime <= p.CurrentTime {
			result = append(result, p.Markers[i])
		}
	}
	return result
}

func (p *Process) DeleteFinishedMarkers() {
	for i := 0; i < len(p.Markers); i++ {
		if p.Markers[i].NextActivationTime <= p.CurrentTime {
			p.Markers = append(p.Markers[:i], p.Markers[i+1:]...)
			i--
		}
	}
}

func (p *Process) UpdateCurrentTime(currentTime float64) {
	p.CurrentTime = currentTime
}
