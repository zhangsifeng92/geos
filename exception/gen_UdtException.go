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

var UdtExceptionName = reflect.TypeOf(UdtException{}).Name()

type UdtException struct {
	Exception
	Elog log.Messages
}

func NewUdtException(parent Exception, message log.Message) *UdtException {
	return &UdtException{parent, log.Messages{message}}
}

func (e UdtException) Code() int64 {
	return UdtErrorCode
}

func (e UdtException) Name() string {
	return UdtExceptionName
}

func (e UdtException) What() string {
	return "UDT error"
}

func (e *UdtException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e UdtException) GetLog() log.Messages {
	return e.Elog
}

func (e UdtException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e UdtException) DetailMessage() string {
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

func (e UdtException) String() string {
	return e.DetailMessage()
}

func (e UdtException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: UdtErrorCode,
		Name: UdtExceptionName,
		What: "UDT error",
	}

	return json.Marshal(except)
}

func (e UdtException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*UdtException):
		callback(&e)
		return true
	case func(UdtException):
		callback(e)
		return true
	default:
		return false
	}
}
