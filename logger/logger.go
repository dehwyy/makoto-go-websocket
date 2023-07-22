package logger

import logg "log"

func Log(s string) {
	l := logg.Default()
	l.Println(s)
}
