package alog

import (
	"runtime"
)

func getGID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// goroutine 1234 [...]
	var id int
	for _, b := range buf[:n] {
		if b >= '0' && b <= '9' {
			id = id*10 + int(b-'0')
		} else if id > 0 {
			break
		}
	}
	return id
}
