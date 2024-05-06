package cps1

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type jamma struct {
	log         *log.Log
	jammaMemory memory.Memory[uint32]
	kickMemory  memory.Memory[uint32]
}

func createJamma(log *log.Log, jammaMemory memory.Memory[uint32], kickMemory memory.Memory[uint32]) *jamma {
	return &jamma{
		log:         log,
		jammaMemory: jammaMemory,
		kickMemory:  kickMemory,
	}
}
