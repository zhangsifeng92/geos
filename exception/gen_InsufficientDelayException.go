// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/zhangsifeng92/geos/log"
)

// template type Exception(PARENT,CODE,WHAT)

var InsufficientDelayExceptionName = reflect.TypeOf(InsufficientDelayException{}).Name()

type InsufficientDelayException struct {
	_AuthorizationException
	Elog log.Messages
}

func NewInsufficientDelayException(parent _AuthorizationException, message log.Message) *InsufficientDelayException {
	return &InsufficientDelayException{parent, log.Messages{message}}
}

func (e InsufficientDelayException) Code() int64 {
	return 3090006
}

func (e InsufficientDelayException) Name() string {
	return InsufficientDelayExceptionName
}

func (e InsufficientDelayException) What() string {
	return "Insufficient delay"
}

func (e *InsufficientDelayException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InsufficientDelayException) GetLog() log.Messages {
	return e.Elog
}

func (e InsufficientDelayException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InsufficientDelayException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e InsufficientDelayException) String() string {
	return e.DetailMessage()
}

func (e InsufficientDelayException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3090006,
		Name: InsufficientDelayExceptionName,
		What: "Insufficient delay",
	}

	return json.Marshal(except)
}

func (e InsufficientDelayException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InsufficientDelayException):
		callback(&e)
		return true
	case func(InsufficientDelayException):
		callback(e)
		return true
	default:
		return false
	}
}
