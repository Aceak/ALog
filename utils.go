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

func parseSize(size interface{}) int64 {
	switch v := size.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		var unit int64 = 1
		switch v[len(v)-1] {
		case 'k', 'K':
			unit = 1024
		case 'm', 'M':
			unit = 1024 * 1024
		case 'g', 'G':
			unit = 1024 * 1024 * 1024
		}
		return int64(float64(parseInt(v[:len(v)-1])) * float64(unit))
	default:
		return 0
	}
}

func parseInt(s string) int {
	var id int
	for _, b := range s {
		if b >= '0' && b <= '9' {
			id = id*10 + int(b-'0')
		} else if id > 0 {
			break
		}
	}
	return id
}
