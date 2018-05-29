package common

import (
	"log"
	"runtime"
)

func PrintStack(all bool) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)

	log.Println("[FATAL] catch a panic,stack is: ", string(buf[:n]))
}
