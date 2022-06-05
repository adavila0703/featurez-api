package api

import (
	"log"

	jsoniter "github.com/json-iterator/go"
)

type Logger struct {
	StatusCode   int
	Method       string
	RequestType  string
	RequestData  string
	ResponseData string
}

func (l *Logger) LogInfo() {
	resp, err := jsoniter.Marshal(l)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(resp))
}
