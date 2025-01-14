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

var LeewayDeadlineExceptionName = reflect.TypeOf(LeewayDeadlineException{}).Name()

type LeewayDeadlineException struct {
	_DeadlineException
	Elog log.Messages
}

func NewLeewayDeadlineException(parent _DeadlineException, message log.Message) *LeewayDeadlineException {
	return &LeewayDeadlineException{parent, log.Messages{message}}
}

func (e LeewayDeadlineException) Code() int64 {
	return 3081001
}

func (e LeewayDeadlineException) Name() string {
	return LeewayDeadlineExceptionName
}

func (e LeewayDeadlineException) What() string {
	return "Transaction reached the deadline set due to leeway on account CPU limits"
}

func (e *LeewayDeadlineException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e LeewayDeadlineException) GetLog() log.Messages {
	return e.Elog
}

func (e LeewayDeadlineException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e LeewayDeadlineException) DetailMessage() string {
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

func (e LeewayDeadlineException) String() string {
	return e.DetailMessage()
}

func (e LeewayDeadlineException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3081001,
		Name: LeewayDeadlineExceptionName,
		What: "Transaction reached the deadline set due to leeway on account CPU limits",
	}

	return json.Marshal(except)
}

func (e LeewayDeadlineException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*LeewayDeadlineException):
		callback(&e)
		return true
	case func(LeewayDeadlineException):
		callback(e)
		return true
	default:
		return false
	}
}
