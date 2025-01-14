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

var KeyNotFoundExceptionName = reflect.TypeOf(KeyNotFoundException{}).Name()

type KeyNotFoundException struct {
	Exception
	Elog log.Messages
}

func NewKeyNotFoundException(parent Exception, message log.Message) *KeyNotFoundException {
	return &KeyNotFoundException{parent, log.Messages{message}}
}

func (e KeyNotFoundException) Code() int64 {
	return KeyNotFoundExceptionCode
}

func (e KeyNotFoundException) Name() string {
	return KeyNotFoundExceptionName
}

func (e KeyNotFoundException) What() string {
	return "Key Not Found"
}

func (e *KeyNotFoundException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e KeyNotFoundException) GetLog() log.Messages {
	return e.Elog
}

func (e KeyNotFoundException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e KeyNotFoundException) DetailMessage() string {
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

func (e KeyNotFoundException) String() string {
	return e.DetailMessage()
}

func (e KeyNotFoundException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: KeyNotFoundExceptionCode,
		Name: KeyNotFoundExceptionName,
		What: "Key Not Found",
	}

	return json.Marshal(except)
}

func (e KeyNotFoundException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*KeyNotFoundException):
		callback(&e)
		return true
	case func(KeyNotFoundException):
		callback(e)
		return true
	default:
		return false
	}
}
