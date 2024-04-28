package m68000

import "github.com/f1gopher/gbpixellib/log"

type M68000 struct {
	log *log.Log
}

func CreateM68000(log *log.Log) *M68000 {
	return &M68000{log: log}
}
