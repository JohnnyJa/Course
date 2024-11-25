package main

import (
	"Course/Interface"
	"Course/Marker"
	"Course/Models"
	"Course/Process"
	"Course/Queue"
	"Course/System"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"math/rand"
)

var failure = 0
var success = 0

var create = Process.NewCreate()
var k1 = Process.NewProcess()
var k2 = Process.NewProcess()
var k3 = Process.NewProcess()
var k4 = Process.NewProcess()
var k5 = Process.NewProcess()
var k6 = Process.NewProcess()

var q1 = Queue.NewQueue()
var q2 = Queue.NewQueue()

var server = System.NewSystem("server", 1000)
var db = System.NewSystem("db", 33)
var model = Models.NewModel(5000, []Interface.IProcess{create, k1, k2, k3, k4, k5, k6})

func main() {
	//stats

	create.GetDelay = func() float64 { return rand.NormFloat64()*0.1 + 0.33 }

	q1.Name = "q1"
	q2.Name = "q2"

	k1.Name = "k1"
	k1.GetDelay = func() float64 { return rand.NormFloat64()*0.5 + 1 }

	k2.Name = "k2"
	k2.GetDelay = func() float64 {
		return rand.NormFloat64()*0.5 + 2
	}

	k3.Name = "k3"
	k3.GetDelay = func() float64 {
		return rand.NormFloat64()*0.5 + 10
	}

	k4.Name = "k4"
	k4.GetDelay = func() float64 {
		return rand.NormFloat64()*300 + 1000
	}

	k5.Name = "k5"
	k5.GetDelay = func() float64 {
		return rand.NormFloat64()*0.5 + 1
	}

	k6.Name = "k6"
	k6.GetDelay = func() float64 {
		return rand.NormFloat64()*0.5 + 10
	}

	create.EndProcess = func() {
		log.Default().Printf("Create process started at %f", create.CurrentTime)
		marker := create.GenerateNewMarker()

		marker.ProcessorsId[0] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: k1.CurrentTime}

		if pr, ok := server.TryStartAtFreeProcessor(); ok {
			marker.ProcessorsId[0].ProcessorId = pr
			k1.TakeMarker(marker)
		} else {

			q1.Push(marker)
			CountMeanTimeByAddingToQueue(create.CurrentTime, "q1")
			log.Default().Printf("Marker %d with type %d pushed to q1. It is %d", marker.Id, marker.Type, q1.Size())

		}
	}

	k1.EndProcess = func() {

		markers := k1.GetFinishedMarkers()
		k1.DeleteFinishedMarkers()

		for _, marker := range markers {

			marker.ProcessorsId[1] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: k1.CurrentTime}
			if pr, ok := db.TryStartAtFreeProcessor(); ok {
				marker.ProcessorsId[1].ProcessorId = pr
				k2.TakeMarker(marker)
			} else {
				q2.Push(marker)
				CountMeanTimeByAddingToQueue(k1.CurrentTime, "q2")
				log.Default().Printf("Marker %d with type %d pushed to q2. It is %d", marker.Id, marker.Type, q2.Size())
			}
		}
	}

	k2.EndProcess = func() {

		markers := k2.GetFinishedMarkers()
		k2.DeleteFinishedMarkers()

		for _, marker := range markers {
			db.FinishProcess(marker.ProcessorsId[1].ProcessorId)

			CountMeanTimeBySystem(marker.ProcessorsId[1].StartTime, k2.CurrentTime, "db")

			marker.ProcessorsId[1] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: -1}
			k3.TakeMarker(marker)
		}

		TryTakeFromQueue(q2, db, k2, k6)
	}

	k3.EndProcess = func() {

		markers := k3.GetFinishedMarkers()
		k3.DeleteFinishedMarkers()

		for _, marker := range markers {
			server.FinishProcess(marker.ProcessorsId[0].ProcessorId)

			CountMeanTimeInNetwork(marker.TimeStart, k6.CurrentTime)
			CountMeanTimeBySystem(marker.ProcessorsId[0].StartTime, k3.CurrentTime, "server")

			marker.ProcessorsId[0] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: -1}
			marker.MaxLifetime = math.MaxFloat64
			k4.TakeMarker(marker)
		}
		TryTakeFromQueue(q1, server, k1, k5)

	}

	k4.EndProcess = func() {
		markers := k4.GetFinishedMarkers()
		k4.DeleteFinishedMarkers()

		for _, marker := range markers {
			marker.Type = 1
			marker.ProcessorsId[0] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: k5.CurrentTime}
			if pr, ok := server.TryStartAtFreeProcessor(); ok {
				marker.ProcessorsId[0].ProcessorId = pr
				k5.TakeMarker(marker)
			} else {
				q1.Push(marker)
				CountMeanTimeByAddingToQueue(k4.CurrentTime, "q1")
				log.Default().Printf("Marker %d with type %d pushed to q1. It is %d", marker.Id, marker.Type, q1.Size())
			}
		}
	}

	k5.EndProcess = func() {

		markers := k5.GetFinishedMarkers()
		k5.DeleteFinishedMarkers()

		for _, marker := range markers {
			marker.ProcessorsId[1] = Marker.ProcessorInfo{ProcessorId: -1, StartTime: k6.CurrentTime}
			if pr, ok := db.TryStartAtFreeProcessor(); ok {
				marker.ProcessorsId[1].ProcessorId = pr
				k6.TakeMarker(marker)
			} else {
				q2.Push(marker)
				CountMeanTimeByAddingToQueue(k5.CurrentTime, "q2")
				log.Default().Printf("Marker %d with type %d pushed to q2. It is %d.", marker.Id, marker.Type, q2.Size())
			}
		}
	}

	k6.EndProcess = func() {

		markers := k6.GetFinishedMarkers()
		k6.DeleteFinishedMarkers()

		for _, marker := range markers {
			server.FinishProcess(marker.ProcessorsId[0].ProcessorId)
			db.FinishProcess(marker.ProcessorsId[1].ProcessorId)

			success++
		}

		TryTakeFromQueue(q1, server, k1, k5)
		TryTakeFromQueue(q2, db, k2, k6)
	}

	model.FinishCycle = func() {
		CheckFailureInProcesses(k1, server, db)
		CheckFailureInProcesses(k2, server, db)
		CheckFailureInProcesses(k3, server, db)
		CheckFailureInProcesses(k4, server, db)
		CheckFailureInProcesses(k5, server, db)
		CheckFailureInProcesses(k6, server, db)
		CheckFailureInQueue(q1, model.CurrentTime, server, db)
		CheckFailureInQueue(q2, model.CurrentTime, server, db)

		meanFailure = append(meanFailure, float64(failure-prevFailure)*math.Abs(model.CurrentTime-prevFailureTime))
		failureTime = append(failureTime, model.CurrentTime)

		prevFailure = failure
		prevFailureTime = model.CurrentTime

		TryTakeFromQueue(q1, server, k1, k5)
		TryTakeFromQueue(q2, db, k2, k6)

		CountMeanQueueSize(q1, model.CurrentTime)
		CountMeanQueueSize(q2, model.CurrentTime)
		previousQueueTime = model.CurrentTime

	}

	model.Simulate()

	fmt.Printf("Success: %d\n", success)
	fmt.Printf("Failure: %d\n", failure)
	fmt.Printf("Mean time by server: %f\n", meanTimeBySystem["server"]/float64(processedBySystem))
	fmt.Printf("Mean time by db: %f\n", meanTimeBySystem["db"]/float64(processedByDB))
	fmt.Printf("Mean time in network: %f\n", meanTimeInNetwork/float64(processedInNetwork))
	fmt.Printf("Mean queue size q1: %f\n", meanQueueSize["q1"]/model.CurrentTime)
	fmt.Printf("Mean queue size q2: %f\n", meanQueueSize["q2"]/model.CurrentTime)
	fmt.Printf("Mean time by adding to q1: %f\n", meanTimeByAddingToQueue["q1"]/float64(addedToQ1))
	fmt.Printf("Mean time by adding to q2: %f\n", meanTimeByAddingToQueue["q2"]/float64(addedToQ2))
	s := 0.0
	for _, f := range meanFailure {
		s += f
	}

	fmt.Printf("Mean failure: %f\n", s/model.CurrentTime)

	makePlot()
}

func TryTakeFromQueue(q *Queue.Queue, system *System.MySystem, kForType0 *Process.Process, kForType1 *Process.Process) {
	for q.Size() > 0 {
		if pr, ok := system.TryStartAtFreeProcessor(); ok {
			m := q.Pop()
			m.ProcessorsId[q.Id].ProcessorId = pr
			if m.Type == 0 {
				kForType0.TakeMarker(m)
				log.Default().Printf("%s get marker %d from %s", kForType0.Name, m.Id, q.Name)
			} else {
				kForType1.TakeMarker(m)
				log.Default().Printf("%s get marker %d from %s", kForType1.Name, m.Id, q.Name)
			}
		} else {
			break
		}
	}
}

var meanTimeBySystem = map[string]float64{"server": 0, "db": 0}
var processedBySystem int
var processedByDB int

func CountMeanTimeBySystem(started float64, finished float64, system string) {
	if system == "server" {
		meanTimeBySystem["server"] += finished - started
		processedBySystem++
	} else {
		meanTimeBySystem["db"] += finished - started
		processedByDB++
	}
}

var meanTimeByAddingToQueue = map[string]float64{"q1": 0, "q2": 0}
var previousQ1Time = 0.0
var previousQ2Time = 0.0

var addedToQ1 = 0
var addedToQ2 = 0

func CountMeanTimeByAddingToQueue(currentTime float64, queue string) {
	if queue == "q1" {
		meanTimeByAddingToQueue["q1"] += currentTime - previousQ1Time
		addedToQ1++
		previousQ1Time = currentTime
	} else {
		meanTimeByAddingToQueue["q2"] += currentTime - previousQ2Time
		addedToQ2++
		previousQ2Time = currentTime
	}
}

var meanTimeInNetwork = 0.0
var processedInNetwork = 0

func CountMeanTimeInNetwork(started float64, finished float64) {
	meanTimeInNetwork += finished - started
	processedInNetwork++
}

var meanQueueSize = map[string]float64{"q1": 0, "q2": 0}

var previousQueueTime = 0.0

func CountMeanQueueSize(q *Queue.Queue, currentTime float64) {
	meanQueueSize[q.Name] += float64(q.Size()) * (currentTime - previousQueueTime)
}

var meanFailure []float64
var failureTime []float64
var prevFailureTime = 0.0
var prevFailure = 0

func CheckFailureInProcesses(p *Process.Process, systems ...*System.MySystem) {
	for i, m := range p.Markers {
		if m.MaxLifetime < p.CurrentTime {
			if len(p.Markers) == 1 {
				p.Markers = []*Marker.Marker{}
			} else {
				p.Markers = append(p.Markers[:i], p.Markers[i+1:]...)
			}
			i--

			log.Default().Printf("Marker %d failed", m.Id)

			for i, system := range systems {
				if m.ProcessorsId[i].ProcessorId != -1 {
					system.FinishProcess(m.ProcessorsId[i].ProcessorId)
				}
			}

			failure++
		}
	}

}

func CheckFailureInQueue(q *Queue.Queue, currentTime float64, systems ...*System.MySystem) {
	for i, m := range q.Elements {
		if m.MaxLifetime < currentTime {
			q.Elements = append(q.Elements[:i], q.Elements[i+1:]...)
			i--
			log.Default().Printf("Marker %d failed", m.Id)

			for j, pr := range m.ProcessorsId {
				if pr.ProcessorId != -1 {
					systems[j].FinishProcess(pr.ProcessorId)
				}
			}
			failure++
		}
	}
}

func makePlot() {
	p := plot.New()
	p.Title.Text = "Time vs Value"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Value"

	// Створення точок для графіка
	pts := make(plotter.XYs, len(failureTime))
	for i := range failureTime {
		pts[i].X = failureTime[i]
		pts[i].Y = meanFailure[i]
	}

	// Додавання лінії до графіка
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	// Збереження графіка у файл
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "time_vs_value.png"); err != nil {
		panic(err)
	}
}
