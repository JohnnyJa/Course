package Marker

var id = 0

type ProcessorInfo struct {
	ProcessorId int
	StartTime   float64
}

type Marker struct {
	Id                 int
	TimeStart          float64
	TimeFinish         float64
	Type               int
	ProcessorsId       [2]ProcessorInfo
	NextActivationTime float64
	MaxLifetime        float64
}

func NewMarker(timeStart float64) *Marker {
	id++
	return &Marker{
		Id:           id,
		TimeStart:    timeStart,
		TimeFinish:   -1,
		ProcessorsId: [2]ProcessorInfo{{-1, -1}, {-1, -1}},
	}
}

func (m *Marker) SetTimeStart(timeStart float64) {
	m.TimeStart = timeStart
}

func NewMarKerWithType(timeStart float64, Type int) *Marker {
	return &Marker{
		TimeStart: timeStart,
		Type:      Type,
	}
}
