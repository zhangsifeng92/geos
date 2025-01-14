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

var AbiCircularDefExceptionName = reflect.TypeOf(AbiCircularDefException{}).Name()

type AbiCircularDefException struct {
	_AbiException
	Elog log.Messages
}

func NewAbiCircularDefException(parent _AbiException, message log.Message) *AbiCircularDefException {
	return &AbiCircularDefException{parent, log.Messages{message}}
}

func (e AbiCircularDefException) Code() int64 {
	return 3150012
}

func (e AbiCircularDefException) Name() string {
	return AbiCircularDefExceptionName
}

func (e AbiCircularDefException) What() string {
	return "Circular definition is detected in the ABI"
}

func (e *AbiCircularDefException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e AbiCircularDefException) GetLog() log.Messages {
	return e.Elog
}

func (e AbiCircularDefException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e AbiCircularDefException) DetailMessage() string {
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

func (e AbiCircularDefException) String() string {
	return e.DetailMessage()
}

func (e AbiCircularDefException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150012,
		Name: AbiCircularDefExceptionName,
		What: "Circular definition is detected in the ABI",
	}

	return json.Marshal(except)
}

func (e AbiCircularDefException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*AbiCircularDefException):
		callback(&e)
		return true
	case func(AbiCircularDefException):
		callback(e)
		return true
	default:
		return false
	}
}
