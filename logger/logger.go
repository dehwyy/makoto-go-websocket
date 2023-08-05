package logger

import "log"

func Log(s string) {
	l := log.Default()
	l.Println(s)
}
