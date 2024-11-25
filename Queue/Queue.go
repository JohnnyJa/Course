package Queue

import (
	"Course/Marker"
)

var id = -1

type Queue struct {
	Id       int
	Elements []*Marker.Marker
	Name     string
}

func NewQueue() *Queue {
	id++
	return &Queue{Id: id}
}

func (q *Queue) Push(element *Marker.Marker) {
	q.Elements = append(q.Elements, element)
}

func (q *Queue) Tail() *Marker.Marker {
	if len(q.Elements) == 0 {
		return nil
	}
	return q.Elements[len(q.Elements)-1]
}

func (q *Queue) Pop() *Marker.Marker {
	if len(q.Elements) == 0 {
		return nil
	}
	element := q.Elements[0]
	q.Elements = q.Elements[1:]
	return element
}

func (q *Queue) PopBack() *Marker.Marker {
	if len(q.Elements) == 0 {
		return nil
	}

	element := q.Elements[len(q.Elements)-1]
	q.Elements = q.Elements[:len(q.Elements)-1]
	return element
}

func (q *Queue) Size() int {
	return len(q.Elements)
}

func (q *Queue) Head() *Marker.Marker {
	if len(q.Elements) == 0 {
		return nil
	}
	return q.Elements[0]
}
