package Process

import "Course/Marker"

type Create struct {
	Process
}

func NewCreate() *Create {
	return &Create{}
}

func (c *Create) GenerateNewMarker() *Marker.Marker {
	marker := Marker.NewMarker(c.CurrentTime)
	marker.Type = 0
	marker.MaxLifetime = c.CurrentTime + 300
	return marker
}

func (c *Create) GetActivationTime() float64 {
	return c.nextActivationTime
}

func (c *Create) RunToCurrentTime() {
	if c.nextActivationTime <= c.CurrentTime {
		c.EndProcess()
		c.nextActivationTime = c.GetDelay() + c.CurrentTime
	}
}
