// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var ProducerScheduleExceptionName = reflect.TypeOf(ProducerScheduleException{}).Name()

type ProducerScheduleException struct {
	_ProducerException
	Elog log.Messages
}

func NewProducerScheduleException(parent _ProducerException, message log.Message) *ProducerScheduleException {
	return &ProducerScheduleException{parent, log.Messages{message}}
}

func (e ProducerScheduleException) Code() int64 {
	return 3170004
}

func (e ProducerScheduleException) Name() string {
	return ProducerScheduleExceptionName
}

func (e ProducerScheduleException) What() string {
	return "Producer schedule exception"
}

func (e *ProducerScheduleException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ProducerScheduleException) GetLog() log.Messages {
	return e.Elog
}

func (e ProducerScheduleException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ProducerScheduleException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("] ")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e ProducerScheduleException) String() string {
	return e.DetailMessage()
}

func (e ProducerScheduleException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3170004,
		Name: ProducerScheduleExceptionName,
		What: "Producer schedule exception",
	}

	return json.Marshal(except)
}

func (e ProducerScheduleException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ProducerScheduleException):
		callback(&e)
		return true
	case func(ProducerScheduleException):
		callback(e)
		return true
	default:
		return false
	}
}
