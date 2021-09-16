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

var InvalidActionClauseExceptionName = reflect.TypeOf(InvalidActionClauseException{}).Name()

type InvalidActionClauseException struct {
	_AbiException
	Elog log.Messages
}

func NewInvalidActionClauseException(parent _AbiException, message log.Message) *InvalidActionClauseException {
	return &InvalidActionClauseException{parent, log.Messages{message}}
}

func (e InvalidActionClauseException) Code() int64 {
	return 3150003
}

func (e InvalidActionClauseException) Name() string {
	return InvalidActionClauseExceptionName
}

func (e InvalidActionClauseException) What() string {
	return "Invalid Ricardian Action"
}

func (e *InvalidActionClauseException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidActionClauseException) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidActionClauseException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InvalidActionClauseException) DetailMessage() string {
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

func (e InvalidActionClauseException) String() string {
	return e.DetailMessage()
}

func (e InvalidActionClauseException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150003,
		Name: InvalidActionClauseExceptionName,
		What: "Invalid Ricardian Action",
	}

	return json.Marshal(except)
}

func (e InvalidActionClauseException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidActionClauseException):
		callback(&e)
		return true
	case func(InvalidActionClauseException):
		callback(e)
		return true
	default:
		return false
	}
}
