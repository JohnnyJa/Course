package System

import "log"

type MySystem struct {
	Name       string
	processors []bool
}

func NewSystem(name string, numOfProcessors int) *MySystem {

	return &MySystem{
		Name:       name,
		processors: make([]bool, numOfProcessors),
	}
}

func (s *MySystem) TryStartAtFreeProcessor() (int, bool) {

	for i, processor := range s.processors {
		if !processor {
			s.processors[i] = true
			log.Printf("Processor %d in %s is busy now", i, s.Name)
			return i, true
		}
	}
	return -1, false
}

func (s *MySystem) FinishProcess(processorIndex int) {
	s.processors[processorIndex] = false
	log.Printf("Processor %d in %s is free", processorIndex, s.Name)
}
