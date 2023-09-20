package errorHandle

import (
	"log"
)

const (
	pattern = "\npath: %s/%s\nfunction: %s\nerror: %s\n"
)

func Commit(path, file, function string, err error) {
	log.Printf(pattern, file, path, function, err)
}

func Fatal(path, file, function string, err error) {
	log.Fatalf(pattern, file, path, function, err)
}
