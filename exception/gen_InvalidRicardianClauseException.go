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

var InvalidRicardianClauseExceptionName = reflect.TypeOf(InvalidRicardianClauseException{}).Name()

type InvalidRicardianClauseException struct {
	_AbiException
	Elog log.Messages
}

func NewInvalidRicardianClauseException(parent _AbiException, message log.Message) *InvalidRicardianClauseException {
	return &InvalidRicardianClauseException{parent, log.Messages{message}}
}

func (e InvalidRicardianClauseException) Code() int64 {
	return 3150002
}

func (e InvalidRicardianClauseException) Name() string {
	return InvalidRicardianClauseExceptionName
}

func (e InvalidRicardianClauseException) What() string {
	return "Invalid Ricardian Clause"
}

func (e *InvalidRicardianClauseException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidRicardianClauseException) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidRicardianClauseException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InvalidRicardianClauseException) DetailMessage() string {
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

func (e InvalidRicardianClauseException) String() string {
	return e.DetailMessage()
}

func (e InvalidRicardianClauseException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150002,
		Name: InvalidRicardianClauseExceptionName,
		What: "Invalid Ricardian Clause",
	}

	return json.Marshal(except)
}

func (e InvalidRicardianClauseException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidRicardianClauseException):
		callback(&e)
		return true
	case func(InvalidRicardianClauseException):
		callback(e)
		return true
	default:
		return false
	}
}
